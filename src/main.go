package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/board"
	"github.com/HyperloopUPV-H8/Backend-H8/connection_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter"
	"github.com/HyperloopUPV-H8/Backend-H8/log_handle"
	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer"
	message_transfer_models "github.com/HyperloopUPV-H8/Backend-H8/message_transfer/models"
	"github.com/HyperloopUPV-H8/Backend-H8/order_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/server"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle"
	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/HyperloopUPV-H8/Backend-H8/websocket_broker"
	"github.com/gorilla/mux"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/browser"
	trace "github.com/rs/zerolog/log"
)

var traceLevel = flag.String("trace", "info", "set the trace level (\"fatal\", \"error\", \"warn\", \"info\", \"debug\", \"trace\")")
var traceFile = flag.String("log", "trace.json", "set the trace log file")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	traceFile := initTrace(*traceLevel, *traceFile)
	defer traceFile.Close()

	config := getConfig()

	excelAdapter := excel_adapter.NewExcelAdapter(config.Excel)

	vehicleBuilder := vehicle.NewBuilder(config.Vehicle)
	podData := vehicle_models.NewPodData()
	orderData := vehicle_models.NewOrderData()

	excelAdapter.Update(vehicleBuilder, podData, orderData)

	vehicle := vehicleBuilder.Build()

	vehicleOutput := make(chan vehicle_models.Update)
	go vehicle.Listen(vehicleOutput)

	boardMux := board.NewMux(board.WithInput(vehicleOutput), board.WithOutput(vehicle.SendOrder))

	updateChan := make(chan vehicle_models.Update)
	go boardMux.Listen(updateChan)

	// Communication with front-end
	websocketBroker := websocket_broker.Get()

	connection_transfer.SetConfig(config.Connections)
	connectionTransfer := connection_transfer.Get()

	data_transfer.SetConfig(config.DataTransfer)
	dataTransfer := data_transfer.Get()

	log_handle.SetConfig(config.Logger)
	logger := log_handle.Get()

	message_transfer.SetConfig(config.Messages)
	messageTransfer := message_transfer.Get()

	orderTransfer, orderChannel := order_transfer.Get()

	websocketBroker.RegisterHandle(connectionTransfer, config.Connections.UpdateTopic)
	websocketBroker.RegisterHandle(dataTransfer)
	websocketBroker.RegisterHandle(logger, config.Logger.Topics.Enable, config.Logger.Topics.State)
	websocketBroker.RegisterHandle(messageTransfer)
	websocketBroker.RegisterHandle(orderTransfer, config.Orders.SendTopic)

	vehicle.OnConnectionChange(connectionTransfer.Update)

	idToType := getIdToType(podData)
	go func() {
		for update := range updateChan {
			logger.Update(update)
			if idToType[update.ID] == "data" {
				dataTransfer.Update(update)
			} else if msg, err := message_transfer_models.MessageFromUpdate(update); err == nil {
				messageTransfer.SendMessage(msg)
			}
		}
	}()

	go func() {
		for order := range orderChannel {
			if err := boardMux.Request(order); err != nil {
				trace.Error().Stack().Err(err).Msg("")
			}
		}
	}()

	httpServer := server.New(mux.NewRouter())

	httpServer.ServeData("/backend/"+config.Server.Endpoints.PodData, podData)
	httpServer.ServeData("/backend/"+config.Server.Endpoints.OrderData, orderData)

	httpServer.HandleFunc("/backend/"+config.Server.Endpoints.Websocket, websocketBroker.HandleConn)

	path, _ := os.Getwd()
	httpServer.FileServer(config.Server.Endpoints.FileServer, filepath.Join(path, config.Server.FileServerPath))

	go httpServer.ListenAndServe(config.Server.Address)
	browser.OpenURL(fmt.Sprintf("http://%s", config.Server.Address))

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

func getIdToType(podData *vehicle_models.PodData) map[uint16]string {
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

func getConfig() Config {
	configFile, fileErr := os.ReadFile("./config.toml")

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
