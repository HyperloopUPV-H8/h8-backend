package main

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"

	blcuPackage "github.com/HyperloopUPV-H8/Backend-H8/blcu"
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
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle"
	vehiclePackage "github.com/HyperloopUPV-H8/Backend-H8/vehicle"
	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/HyperloopUPV-H8/Backend-H8/websocket_broker"
	"github.com/gorilla/mux"
	"github.com/pelletier/go-toml/v2"
	trace "github.com/rs/zerolog/log"
)

var traceLevel = flag.String("trace", "info", "set the trace level (\"fatal\", \"error\", \"warn\", \"info\", \"debug\", \"trace\")")
var traceFile = flag.String("log", "trace.json", "set the trace log file")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	traceFile := initTrace(*traceLevel, *traceFile)
	defer traceFile.Close()

	config := getConfig("./config.toml")

	excelAdapter := excel_adapter.New(config.Excel)
	boards := excelAdapter.GetBoards()
	globalInfo := excelAdapter.GetGlobalInfo()

	connectionTransfer := connection_transfer.New(config.Connections)

	podData := vehicle_models.NewPodData(boards)
	orderData := vehicle_models.NewOrderData(boards)
	blcu := blcuPackage.NewBLCU(globalInfo, config.BLCU)
	uploadableBoards := blcuPackage.GetUploadableBoards(globalInfo, config.Excel.Parse.Global.BLCUAddressKey)

	vehicle := vehicle.New(vehicle.VehicleConstructorArgs{
		Config:             config.Vehicle,
		Boards:             boards,
		GlobalInfo:         globalInfo,
		OnConnectionChange: connectionTransfer.Update,
	})
	vehicleUpdates := make(chan vehicle_models.PacketUpdate, 100)
	vehicleProtections := make(chan vehicle_models.Protection)
	// vehicleOrders := make(chan packet.Packet)
	go vehicle.Listen(vehicleUpdates, vehicleProtections)

	dataTransfer := data_transfer.New(config.DataTransfer)
	go dataTransfer.Run()

	messageTransfer := message_transfer.New(config.Messages)
	orderTransfer, orderChannel := order_transfer.New()

	packetLogger := packet_logger.NewPacketLogger(boards, config.PacketLogger)
	orderLogger := order_logger.NewOrderLogger(boards, config.OrderLogger)
	protectionLogger := protection_logger.NewProtectionLogger(config.ProtectionLogger, config.Vehicle.Protections)

	loggers := map[string]logger_handler.Logger{
		"packet":     &packetLogger,
		"order":      &orderLogger,
		"protection": &protectionLogger,
	}

	loggerHandler := logger_handler.NewLoggerHandler(loggers, config.LoggerHandler)

	websocketBroker := websocket_broker.New()
	defer websocketBroker.Close()

	websocketBroker.RegisterHandle(&blcu, config.BLCU.Topics.Upload, config.BLCU.Topics.Download)
	websocketBroker.RegisterHandle(&connectionTransfer, config.Connections.UpdateTopic)
	websocketBroker.RegisterHandle(&dataTransfer)
	websocketBroker.RegisterHandle(&loggerHandler, config.LoggerHandler.Topics.Enable, config.LoggerHandler.Topics.State)
	websocketBroker.RegisterHandle(&messageTransfer)
	websocketBroker.RegisterHandle(&orderTransfer, config.Orders.SendTopic)

	go startPacketUpdateRoutine(vehicleUpdates, &dataTransfer, &loggerHandler)
	go startProtectionsRoutine(vehicleProtections, &messageTransfer, &loggerHandler)
	go startOrderRoutine(orderChannel, &vehicle, &loggerHandler)

	// go func() {
	// 	for packet := range vehicleOrders {
	// 		logger.Update(packet)
	// 	}
	// }()

	go func() {
		for id := range websocketBroker.CloseChan {
			loggerHandler.NotifyDisconnect(id)
		}
	}()

	httpServer := server.New(mux.NewRouter())

	httpServer.ServeData("/backend"+config.Server.Endpoints.PodData, podData)
	httpServer.ServeData("/backend"+config.Server.Endpoints.OrderData, orderData)
	httpServer.ServeData("/backend"+config.Server.Endpoints.UploadableBoards, uploadableBoards)

	httpServer.HandleFunc(config.Server.Endpoints.Websocket, websocketBroker.HandleConn)

	path, _ := os.Getwd()
	httpServer.FileServer(config.Server.Endpoints.FileServer, filepath.Join(path, config.Server.FileServerPath))

	go httpServer.ListenAndServe(config.Server.Address)

	// browser.OpenURL(fmt.Sprintf("http://%s", config.Server.Address))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
}

func startPacketUpdateRoutine(vehicleUpdates <-chan vehicle_models.PacketUpdate, dataTransfer *data_transfer.DataTransfer, loggerHandler *logger_handler.LoggerHandler) {
	updateFactory := update_factory.NewFactory()

	for packetUpdate := range vehicleUpdates {
		update := updateFactory.NewUpdate(packetUpdate)
		dataTransfer.Update(update)

		loggerHandler.Log(packet_logger.ToLoggablePacket(packetUpdate))

		for id, value := range packetUpdate.Values {
			loggerHandler.Log(packet_logger.ToLoggableValue(id, value, packetUpdate.Metadata.Timestamp))
		}
	}
}

func startProtectionsRoutine(vehicleProtections <-chan vehicle_models.Protection, messageTransfer *message_transfer.MessageTransfer, loggerHandler *logger_handler.LoggerHandler) {
	for protection := range vehicleProtections {
		messageTransfer.SendMessage(protection)
		loggerHandler.Log(protection_logger.LoggableProtection(protection))
	}
}

func startOrderRoutine(orderChannel <-chan vehicle_models.Order, vehicle *vehiclePackage.Vehicle, loggerHandler *logger_handler.LoggerHandler) {
	for ord := range orderChannel {
		err := vehicle.SendOrder(ord)

		if err != nil {
			trace.Error().Any("order", ord).Msg("error sending order")
		}

		loggerHandler.Log(order_logger.LoggableOrder(ord))
	}
}

func getConfig(path string) Config {
	configFile, fileErr := os.ReadFile(path)

	if fileErr != nil {
		trace.Fatal().Stack().Err(fileErr).Msg("error reading config file")
	}

	reader := strings.NewReader(string(configFile))

	var config Config

	// decodeErr := toml.NewDecoder(reader).DisallowUnknownFields().Decode(&config)
	decodeErr := toml.NewDecoder(reader).Decode(&config)

	if decodeErr != nil {
		trace.Fatal().Stack().Err(decodeErr).Msg("error unmarshaling toml file")
	}

	return config
}
