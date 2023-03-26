package main

import (
	"fmt"
	"log"
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
	log_handle_models "github.com/HyperloopUPV-H8/Backend-H8/log_handle/models"
	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/order_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/server"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle"
	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/HyperloopUPV-H8/Backend-H8/websocket_handle"
	"github.com/HyperloopUPV-H8/Backend-H8/websocket_handle/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	godotenv.Load(".env")

	document := excel_adapter.FetchDocument(os.Getenv("EXCEL_ID"), os.Getenv("EXCEL_PATH"), os.Getenv("EXCEL_NAME"))

	vehicleBuilder := vehicle.NewBuilder()
	podData := vehicle_models.NewPodData()
	orderData := vehicle_models.NewOrderData()

	excel_adapter.Update(document, vehicleBuilder, podData, orderData)

	vehicle := vehicleBuilder.Build()

	vehicleOutput := make(chan vehicle_models.Update)
	go vehicle.Listen(vehicleOutput)

	boardMux := board.NewMux(board.WithInput(vehicleOutput), board.WithOutput(vehicle.SendOrder))
	log.Println("New Mux")

	updateChan := make(chan vehicle_models.Update)
	go boardMux.Listen(updateChan)
	log.Println("Mux Listen")

	idToType := getIdToType(podData)

	connectionTransfer, connectionChannel := connection_transfer.New()
	vehicle.OnConnectionChange(connectionTransfer.Update)

	dataTransfer, dataTransferChannel := data_transfer.New(getFPS(30))

	messageTransfer, messageChannel := message_transfer.New()

	orderChannel := make(chan vehicle_models.Order, 100)
	_, ordChannel := order_transfer.New(orderChannel)

	logger, loggerChannel := log_handle.NewLogger(log_handle_models.Config{
		DumpSize: 7000,
		RowSize:  20,
		BasePath: os.Getenv("LOG_PATH"),
		Updates:  make(chan map[string]any, 10000),
		Autosave: time.NewTicker(time.Minute),
	})

	go func(msgChannel chan models.MessageTarget) {
		for update := range updateChan {
			logger.Update(update.Fields)
			if idToType[update.ID] == "data" {
				dataTransfer.Update(update)
			} else {
				messageTransfer.Broadcast(update)
			}
		}
	}(messageChannel)

	go func() {
		for order := range orderChannel {
			log.Println(order)
			if err := boardMux.Request(order); err != nil {
				log.Printf("request failed: %s\n", err)
			}
		}
	}()

	httpServer := server.Server{Router: mux.NewRouter()}

	httpServer.ServeData("/backend/"+os.Getenv("POD_DATA_ENDPOINT"), podData)
	httpServer.ServeData("/backend/"+os.Getenv("ORDER_DATA_ENDPOINT"), orderData)

	websocket_handle.RunWSHandle(httpServer.Router, "/backend", map[string]chan models.MessageTarget{
		"podData/update":    dataTransferChannel,
		"message/update":    messageChannel,
		"order/send":        ordChannel,
		"connection/update": connectionChannel,
		"logger":            loggerChannel,
	})

	path, _ := os.Getwd()
	httpServer.FileServer("/", filepath.Join(path, "static"))

	go httpServer.ListenAndServe(os.Getenv("SERVER_ADDR"))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Println("backend running!")
loop:
	for {
		select {
		case <-time.After(time.Second * 10):
			fmt.Println(vehicle.Stats())
		case <-interrupt:
			break loop
		}
	}
}

func getFPS(fps int) time.Duration {
	return time.Duration(int(time.Second) / fps)
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
