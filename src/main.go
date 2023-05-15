package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"

	blcuPackage "github.com/HyperloopUPV-H8/Backend-H8/blcu"
	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/connection_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter"
	"github.com/HyperloopUPV-H8/Backend-H8/logger_handler"
	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/order_logger"
	"github.com/HyperloopUPV-H8/Backend-H8/order_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/packet_logger"
	"github.com/HyperloopUPV-H8/Backend-H8/protection_logger"
	"github.com/HyperloopUPV-H8/Backend-H8/server"
	"github.com/HyperloopUPV-H8/Backend-H8/update_factory"
	"github.com/HyperloopUPV-H8/Backend-H8/value_logger"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle"
	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/HyperloopUPV-H8/Backend-H8/websocket_broker"
	"github.com/fatih/color"
	"github.com/google/gopacket/pcap"
	"github.com/pelletier/go-toml/v2"
	trace "github.com/rs/zerolog/log"
)

var traceLevel = flag.String("trace", "info", "set the trace level (\"fatal\", \"error\", \"warn\", \"info\", \"debug\", \"trace\")")
var traceFile = flag.String("log", "trace.json", "set the trace log file")

func main() {

	traceFile := initTrace(*traceLevel, *traceFile)
	defer traceFile.Close()

	pidPath := path.Join(os.TempDir(), "backendPid")

	createPid(pidPath)
	defer RemovePid(pidPath)

	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	config := getConfig("./config.toml")
	log.Println(config.Test)

	excelAdapter := excel_adapter.New(config.Excel)
	boards := excelAdapter.GetBoards()
	globalInfo := excelAdapter.GetGlobalInfo()

	dev, err := selectDev()
	if err != nil {
		trace.Fatal().Err(err).Msg("Error selecting device")
		panic(err)
	}
	config.Vehicle.Network.Interface = dev.Name

	connectionTransfer := connection_transfer.New(config.Connections)

	podData := vehicle_models.NewPodData(boards)
	orderData := vehicle_models.NewOrderData(boards, config.Excel.Parse.Global.BLCUAddressKey)

	vehicle := vehicle.New(vehicle.VehicleConstructorArgs{
		Config:             config.Vehicle,
		Boards:             boards,
		GlobalInfo:         globalInfo,
		OnConnectionChange: connectionTransfer.Update,
	})

	blcu := blcuPackage.NewBLCU(globalInfo, config.BLCU)
	blcu.SetSendOrder(vehicle.SendOrder)

	vehicleUpdates := make(chan vehicle_models.PacketUpdate, 1)
	vehicleProtections := make(chan vehicle_models.ProtectionMessage)
	vehicleTransmittedOrders := make(chan vehicle_models.PacketUpdate)
	blcuAckChan := make(chan struct{})

	dataTransfer := data_transfer.New(config.DataTransfer)
	go dataTransfer.Run()

	messageTransfer := message_transfer.New(config.Messages)
	orderTransfer, orderChannel := order_transfer.New()

	packetLogger := packet_logger.NewPacketLogger(boards, config.PacketLogger)
	valueLogger := value_logger.NewValueLogger(boards, config.ValueLogger)
	orderLogger := order_logger.NewOrderLogger(boards, config.OrderLogger)
	protectionLogger := protection_logger.NewProtectionLogger(config.Vehicle.Messages.FaultIdKey, config.Vehicle.Messages.WarningIdKey, config.Vehicle.Messages.ErrorIdKey, config.ProtectionLogger)

	loggers := map[string]logger_handler.Logger{
		"packets":     &packetLogger,
		"values":      &valueLogger,
		"orders":      &orderLogger,
		"protections": &protectionLogger,
	}

	loggerHandler := logger_handler.NewLoggerHandler(loggers, config.LoggerHandler)

	websocketBroker := websocket_broker.New()
	defer websocketBroker.Close()

	websocketBroker.RegisterHandle(&blcu, config.BLCU.Topics.Upload, config.BLCU.Topics.Download)
	websocketBroker.RegisterHandle(&connectionTransfer, config.Connections.UpdateTopic)
	websocketBroker.RegisterHandle(&dataTransfer)
	websocketBroker.RegisterHandle(&loggerHandler, config.LoggerHandler.Topics.Enable)
	websocketBroker.RegisterHandle(&messageTransfer)
	websocketBroker.RegisterHandle(&orderTransfer, config.Orders.SendTopic)

	go vehicle.Listen(vehicleUpdates, vehicleTransmittedOrders, vehicleProtections, blcuAckChan)

	go startPacketUpdateRoutine(vehicleUpdates, &dataTransfer, &loggerHandler)
	go startProtectionsRoutine(vehicleProtections, &messageTransfer, &loggerHandler)
	go startOrderRoutine(orderChannel, &vehicle, &loggerHandler)

	go func() {
		for order := range vehicleTransmittedOrders {
			loggable := order_logger.LoggableTransmittedOrder(order)
			loggerHandler.Log(loggable)
		}
	}()

	go func() {
		for range blcuAckChan {
			blcu.NotifyAck()
		}
	}()

	go func() {
		for id := range websocketBroker.CloseChan {
			loggerHandler.NotifyDisconnect(id)
		}
	}()

	uploadableBords := common.Filter(common.Keys(globalInfo.BoardToIP), func(item string) bool {
		return item != config.Excel.Parse.Global.BLCUAddressKey
	})

	endpointData := server.EndpointData{
		PodData:           podData,
		OrderData:         orderData,
		ProgramableBoards: uploadableBords,
	}

	serverHandler, err := server.New(&websocketBroker, endpointData, config.Server)
	if err != nil {
		trace.Fatal().Err(err).Msg("Error creating server")
		panic(err)
	}

	errs := serverHandler.ListenAndServe()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case err := <-errs:
			trace.Error().Err(err).Msg("Error in server")

		case <-interrupt:
			trace.Info().Msg("Shutting down")
			return
		}
	}
}

func createPid(path string) {
	err := WritePid(path)

	if err != nil {
		switch err {
		case ErrProcessRunning:
			trace.Fatal().Err(err).Msg("Backend is already running")
		default:
			trace.Error().Err(err).Msg("pid error")
		}
	}
}

func selectDev() (pcap.Interface, error) {
	devs, err := pcap.FindAllDevs()
	if err != nil {
		return pcap.Interface{}, err
	}

	cyan := color.New(color.FgCyan)

	cyan.Print("select a device: ")
	fmt.Printf("(0-%d)\n", len(devs)-1)
	for i, dev := range devs {
		displayDev(i, dev)
	}

	dev, err := acceptInput(len(devs))
	if err != nil {
		return pcap.Interface{}, err
	}

	return devs[dev], nil
}

func displayDev(i int, dev pcap.Interface) {
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	red.Printf("\t%d", i)
	fmt.Print(": (")
	yellow.Print(dev.Name)
	fmt.Printf(") %s [", dev.Description)
	for _, addr := range dev.Addresses {
		green.Printf("%s", addr.IP)
		fmt.Print(", ")
	}
	fmt.Println("]")
}

func acceptInput(limit int) (int, error) {
	blue := color.New(color.FgBlue)
	red := color.New(color.FgRed)

	for {
		blue.Print(">>> ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return 0, err
		}

		var dev int
		_, err = fmt.Sscanf(input, "%d", &dev)
		if err != nil {
			red.Printf("%s\n\n", err)
			continue
		}

		if dev < 0 || dev >= limit {
			red.Println("invalid device selected\n")
			continue
		} else {
			return dev, nil
		}
	}
}

func getConfig(path string) Config {
	configFile, fileErr := os.ReadFile(path)

	if fileErr != nil {
		trace.Fatal().Stack().Err(fileErr).Msg("error reading config file")
	}

	reader := strings.NewReader(string(configFile))

	var config Config

	// TODO: add strict mode (DisallowUnkownFields)
	decodeErr := toml.NewDecoder(reader).Decode(&config)

	if decodeErr != nil {
		trace.Fatal().Stack().Err(decodeErr).Msg("error unmarshaling toml file")
	}

	return config
}

func startPacketUpdateRoutine(vehicleUpdates <-chan vehicle_models.PacketUpdate, dataTransfer *data_transfer.DataTransfer, loggerHandler *logger_handler.LoggerHandler) {
	updateFactory := update_factory.NewFactory()

	for packetUpdate := range vehicleUpdates {
		update := updateFactory.NewUpdate(packetUpdate)
		dataTransfer.Update(update)

		loggerHandler.Log(packet_logger.ToLoggablePacket(packetUpdate))

		for id, value := range packetUpdate.Values {
			loggerHandler.Log(value_logger.ToLoggableValue(id, value, packetUpdate.Metadata.Timestamp))
		}
	}
}

func startProtectionsRoutine(vehicleProtections <-chan vehicle_models.ProtectionMessage, messageTransfer *message_transfer.MessageTransfer, loggerHandler *logger_handler.LoggerHandler) {
	for protection := range vehicleProtections {
		messageTransfer.SendMessage(protection)
		loggerHandler.Log(protection_logger.LoggableProtection(protection))
	}
}

func startOrderRoutine(orderChannel <-chan vehicle_models.Order, vehicle *vehicle.Vehicle, loggerHandler *logger_handler.LoggerHandler) {
	for ord := range orderChannel {
		err := vehicle.SendOrder(ord)

		if err != nil {
			trace.Error().Any("order", ord).Msg("error sending order")
		}

		loggerHandler.Log(order_logger.LoggableOrder(ord))
	}
}
