package main

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/blcu"
	"github.com/HyperloopUPV-H8/Backend-H8/connection_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter"
	loggerPackage "github.com/HyperloopUPV-H8/Backend-H8/logger"
	message_parser_models "github.com/HyperloopUPV-H8/Backend-H8/message_parser/models"
	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/order_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/server"
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
	blcu := blcu.NewBLCU(globalInfo, config.BLCU)

	vehicle := vehicle.NewVehicle(boards, globalInfo, config.Vehicle, connectionTransfer.Update)
	vehicleUpdates := make(chan vehicle_models.Update)
	vehicleMessages := make(chan interface{})
	go vehicle.Listen(vehicleUpdates, vehicleMessages)

	dataTransfer := data_transfer.New(config.DataTransfer)
	go dataTransfer.Run()

	messageTransfer := message_transfer.New(config.Messages)
	orderTransfer, orderChannel := order_transfer.New()
	logger := loggerPackage.New(config.Logger)
	go logger.Listen()

	// Communication with front-end
	websocketBroker := websocket_broker.New()
	defer websocketBroker.Close()

	websocketBroker.RegisterHandle(&blcu, config.BLCU.Topics.Upload, config.BLCU.Topics.Download)
	websocketBroker.RegisterHandle(&connectionTransfer, config.Connections.UpdateTopic)
	websocketBroker.RegisterHandle(&dataTransfer)
	websocketBroker.RegisterHandle(&logger, config.Logger.Topics.Enable, config.Logger.Topics.State)
	websocketBroker.RegisterHandle(&messageTransfer)
	websocketBroker.RegisterHandle(&orderTransfer, config.Orders.SendTopic)

	idToType := getIdToType(podData)
	go func() {
		for update := range vehicleUpdates {
			logger.Updates <- update
			if idToType[update.ID] == "data" {
				dataTransfer.Update(update)
			}
		}
	}()

	go func() {
		for id := range websocketBroker.CloseCh {
			logger.Enable <- loggerPackage.EnableMsg{
				Client: id,
				Enable: false,
			}
		}
	}()

	go func() {
		for message := range vehicleMessages {
			switch m := message.(type) {
			case message_parser_models.ProtectionMessage:
				err := messageTransfer.SendMessage(m)

				if err != nil {
					trace.Error().Err(err).Stack().Msg("error sending message")
				}
			case message_parser_models.BLCU_ACK:
				blcu.HandleACK()
			}
		}
	}()

	go func() {
		for order := range orderChannel {
			vehicle.SendOrder(order)
		}
	}()

	httpServer := server.New(mux.NewRouter())

	httpServer.ServeData("/backend"+config.Server.Endpoints.PodData, podData)
	httpServer.ServeData("/backend"+config.Server.Endpoints.OrderData, orderData)

	httpServer.HandleFunc(config.Server.Endpoints.Websocket, websocketBroker.HandleConn)

	path, _ := os.Getwd()
	httpServer.FileServer(config.Server.Endpoints.FileServer, filepath.Join(path, config.Server.FileServerPath))

	go httpServer.ListenAndServe(config.Server.Address)

	// browser.OpenURL(fmt.Sprintf("http://%s", config.Server.Address))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

loop:
	for {
		select {
		case <-time.After(time.Second * 10):
			trace.Trace().Any("stats", vehicle.Stats()).Msg("stats")
		case <-interrupt:
			break loop
		}
	}
}

func getIdToType(podData vehicle_models.PodData) map[uint16]string {
	idToType := make(map[uint16]string)
	for _, brd := range podData.Boards {
		for _, pkt := range brd.Packets {
			idToType[pkt.ID] = "data"
		measurements_loop:
			for msr := range pkt.Measurements {
				if msr == "warning" {
					idToType[pkt.ID] = "warning"
					break measurements_loop
				} else if msr == "fault" {
					idToType[pkt.ID] = "fault"
					break measurements_loop
				}
			}
		}
	}
	return idToType
}

func getConfig(path string) Config {
	configFile, fileErr := os.ReadFile(path)

	if fileErr != nil {
		trace.Fatal().Stack().Err(fileErr).Msg("error reading config file")
	}

	var config Config
	unmarshalErr := toml.Unmarshal(configFile, &config)

	if unmarshalErr != nil {
		trace.Fatal().Stack().Err(unmarshalErr).Msg("error unmarshaling toml file")
	}

	return config
}
