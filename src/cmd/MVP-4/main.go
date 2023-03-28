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
	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer"
	message_transfer_models "github.com/HyperloopUPV-H8/Backend-H8/message_transfer/models"
	"github.com/HyperloopUPV-H8/Backend-H8/order_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/server"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle"
	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/HyperloopUPV-H8/Backend-H8/websocket_broker"
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

	updateChan := make(chan vehicle_models.Update)
	go boardMux.Listen(updateChan)

	// Communication with front-end
	websocketBroker := websocket_broker.Get()

	connectionTransfer := connection_transfer.Get()
	dataTransfer := data_transfer.Get()
	logger := log_handle.Get()
	messageTransfer := message_transfer.Get()
	orderTransfer, orderChannel := order_transfer.Get()

	websocketBroker.RegisterHandle(connectionTransfer, "connection/get")
	websocketBroker.RegisterHandle(dataTransfer)
	websocketBroker.RegisterHandle(logger, "logger/enable")
	websocketBroker.RegisterHandle(messageTransfer)
	websocketBroker.RegisterHandle(orderTransfer, "order/send")

	vehicle.OnConnectionChange(connectionTransfer.Update)

	go func() {
		for update := range updateChan {
			logger.Update(update)
			if msg, err := message_transfer_models.MessageFromUpdate(update); err != nil {
				dataTransfer.Update(update)
			} else {
				messageTransfer.SendMessage(msg)
			}
		}
	}()

	go func() {
		for order := range orderChannel {
			if err := boardMux.Request(order); err != nil {
				log.Printf("request failed: %s\n", err)
			}
		}
	}()

	httpServer := server.Server{Router: mux.NewRouter()}

	httpServer.ServeData("/backend/"+os.Getenv("POD_DATA_ENDPOINT"), podData)
	httpServer.ServeData("/backend/"+os.Getenv("ORDER_DATA_ENDPOINT"), orderData)

	httpServer.HandleFunc("/backend", websocketBroker.HandleConn)

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
