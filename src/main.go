package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"

	blcuPackage "github.com/HyperloopUPV-H8/Backend-H8/blcu"
	"github.com/HyperloopUPV-H8/Backend-H8/connection_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter"
	loggerPackage "github.com/HyperloopUPV-H8/Backend-H8/logger"
	data_logger "github.com/HyperloopUPV-H8/Backend-H8/logger/data"
	message_logger "github.com/HyperloopUPV-H8/Backend-H8/logger/message"
	order_logger "github.com/HyperloopUPV-H8/Backend-H8/logger/order"
	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/order_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/data"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/message"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/order"
	"github.com/HyperloopUPV-H8/Backend-H8/server"
	"github.com/HyperloopUPV-H8/Backend-H8/update_factory"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle"
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
	vehicleUpdates := make(chan packet.Packet, 100)
	vehicleMessages := make(chan packet.Packet)
	vehicleOrders := make(chan packet.Packet)
	go vehicle.Listen(vehicleUpdates, vehicleMessages, vehicleOrders)

	dataTransfer := data_transfer.New(config.DataTransfer)
	go dataTransfer.Run()

	messageTransfer := message_transfer.New(config.Messages)
	orderTransfer, orderChannel := order_transfer.New()

	dataLogger, err := data_logger.New(data_logger.Config{
		BasePath: config.Logger.BasePath,
		FileName: config.Logger.DataFileName,
	})

	if err != nil {
		panic(err)
	}

	messageLogger, err := message_logger.New(message_logger.Config{
		BasePath: config.Logger.BasePath,
		FileName: config.Logger.MessageFileName,
	})

	if err != nil {
		panic(err)
	}

	orderLogger, err := order_logger.New(order_logger.Config{
		BasePath: config.Logger.BasePath,
		FileName: config.Logger.OrderFileName,
	})

	if err != nil {
		panic(err)
	}

	logger := loggerPackage.New(map[packet.Kind]loggerPackage.SubLogger{
		packet.Data:    dataLogger,
		packet.Message: messageLogger,
		packet.Order:   orderLogger,
	}, config.Logger)
	defer logger.Close()

	// Communication with front-end
	websocketBroker := websocket_broker.New()
	defer websocketBroker.Close()

	websocketBroker.RegisterHandle(&blcu, config.BLCU.Topics.Upload, config.BLCU.Topics.Download)
	websocketBroker.RegisterHandle(&connectionTransfer, config.Connections.UpdateTopic)
	websocketBroker.RegisterHandle(&dataTransfer)
	websocketBroker.RegisterHandle(&logger, config.Logger.Topics.Enable, config.Logger.Topics.State)
	websocketBroker.RegisterHandle(&messageTransfer)
	websocketBroker.RegisterHandle(&orderTransfer, config.Orders.SendTopic)

	updateFactory := update_factory.NewFactory()
	go func() {
		for packet := range vehicleUpdates {
			payload, ok := packet.Payload.(data.Payload)
			if !ok {
				// TODO: handle error
				continue
			}
			logger.Update(packet)

			update := updateFactory.NewUpdate(packet.Metadata, payload)
			dataTransfer.Update(update)
		}
	}()

	go func() {
		for packet := range vehicleMessages {
			payload, ok := packet.Payload.(message.Payload)
			if !ok {
				// TODO: handle error
				continue
			}
			logger.Update(packet)

			messageTransfer.SendMessage(payload.Data)
		}
	}()

	go func() {
		for packet := range vehicleOrders {
			logger.Update(packet)
		}
	}()

	go func() {
		for id := range websocketBroker.CloseChan {
			logger.NotifyDisconnect(id)
		}
	}()

	go func() {
		for ord := range orderChannel {
			log.Println(ord)
			id, fields := convertOrder(ord)
			values, enabled := unzipFields(fields)
			meta, err := vehicle.SendOrder(id, order.Payload{Values: values, Enabled: enabled})
			if err == nil {
				logger.Update(packet.Packet{Metadata: meta, Payload: order.Payload{
					Values:  values,
					Enabled: enabled,
				}})
			}
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

func getConfig(path string) Config {
	configFile, fileErr := os.ReadFile(path)

	if fileErr != nil {
		trace.Fatal().Stack().Err(fileErr).Msg("error reading config file")
	}

	reader := strings.NewReader(string(configFile))

	var config Config
	decodeErr := toml.NewDecoder(reader).Decode(&config)

	if decodeErr != nil {
		trace.Fatal().Stack().Err(decodeErr).Msg("error unmarshaling toml file")
	}

	return config
}

func unzipFields(fields map[string]vehicle_models.Field) (map[string]packet.Value, map[string]bool) {
	fieldsMap := make(map[string]packet.Value)
	enabledMap := make(map[string]bool)

	for name, field := range fields {
		fieldsMap[name] = field.Value.(packet.Value)
		enabledMap[name] = field.IsEnabled
	}

	return fieldsMap, enabledMap
}

func convertOrder(order vehicle_models.Order) (uint16, map[string]vehicle_models.Field) {
	fields := make(map[string]vehicle_models.Field)
	for name, field := range order.Fields {
		newField := vehicle_models.Field{
			IsEnabled: field.IsEnabled,
		}
		switch value := field.Value.(type) {
		case float64:
			newField.Value = packet.Numeric{Value: value}
		case string:
			newField.Value = packet.Enum{Value: value}
		case bool:
			newField.Value = packet.Boolean{Value: value}
		default:
			log.Printf("name: %s, type: %T\n", name, field.Value)
			continue
		}
		fields[name] = newField
	}

	return order.ID, fields
}
