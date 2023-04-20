package main

import (
	"flag"
	"fmt"
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
	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/order_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/data"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/message"
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
	logger := loggerPackage.New(config.Logger)
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
			trace.Debug().Msgf("Received update: %v", packet)
			payload, ok := packet.Payload.(data.Payload)
			if !ok {
				// TODO: handle error
				continue
			}
			update := updateFactory.NewUpdate(packet.Metadata, payload)
			logger.UpdateData(update)
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
			logger.UpdateMsg(payload.Data.String())
			messageTransfer.SendMessage(payload.Data)
		}
	}()

	go func() {
		for packet := range vehicleOrders {
			fmt.Println(packet)
		}
	}()

	go func() {
		for id := range websocketBroker.CloseChan {
			logger.NotifyDisconnect(id)
		}
	}()

	go func() {
		for order := range orderChannel {
			log.Println(order)
			vehicle.SendOrder(order)
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

	reader := strings.NewReader(string(configFile))

	var config Config
	decodeErr := toml.NewDecoder(reader).Decode(&config)

	if decodeErr != nil {
		trace.Fatal().Stack().Err(decodeErr).Msg("error unmarshaling toml file")
	}

	return config
}
