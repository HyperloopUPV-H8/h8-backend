package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	logger "github.com/HyperloopUPV-H8/Backend-H8/Shared/Logger/infra"
	loggerDto "github.com/HyperloopUPV-H8/Backend-H8/Shared/Logger/infra/dto"
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter"
	excelRetriever "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_retriever"
	packetAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra/dto"
	transportControllerInfra "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/transport_controller/infra"
	server "github.com/HyperloopUPV-H8/Backend-H8/Shared/server/infra"
	serverMappers "github.com/HyperloopUPV-H8/Backend-H8/Shared/server/infra/mappers"
	dataTransferApplication "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/application"
	dataTransfer "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/infra"
	dataTransferPresentation "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/infra/presentation"
	messageTransferApplication "github.com/HyperloopUPV-H8/Backend-H8/message_transfer/application"
	messageTransferMappers "github.com/HyperloopUPV-H8/Backend-H8/message_transfer/infra/mappers"
	messageTransferPresentation "github.com/HyperloopUPV-H8/Backend-H8/message_transfer/infra/presentation"
	orderTransferApplication "github.com/HyperloopUPV-H8/Backend-H8/order_transfer/application"
	orderTransferInfra "github.com/HyperloopUPV-H8/Backend-H8/order_transfer/infra/mappers"
	orderTransferPresentation "github.com/HyperloopUPV-H8/Backend-H8/order_transfer/infra/presentation"
	"github.com/joho/godotenv"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	log.Println("Starting Backend...")

	godotenv.Load(".env")

	document := excelRetriever.GetExcel(os.Getenv("SPREADSHEET_ID"), os.Getenv("SPREADSHEET_NAME"), os.Getenv("SPREADSHEET_PATH"), os.Getenv("CREDENTIALS"))
	boards := excelAdapter.GetBoards(document)

	config := transportControllerInfra.Config{
		Device:        "\\Device\\NPF_Loopback",
		Live:          true,
		TCPConfig:     nil,
		SnifferConfig: nil,
	}

	dataChan := make(chan dto.PacketUpdate)
	messageChan := make(chan dto.PacketUpdate)
	orderChan := make(chan dto.PacketUpdate)

	packetFactory := dataTransfer.NewFactory()

	logger := logger.NewLogger(os.Getenv("LOG_DIR"), time.Second*10)

	server := server.New[
		dataTransferApplication.PacketJSON,
		orderTransferApplication.OrderJSON,
		messageTransferApplication.MessageJSON]()

	go func() {
		for {
			packet := <-dataChan
			json := dataTransferApplication.NewJSON(packetFactory.NewPacket(packet))

		loop:
			for name, measure := range packet.Values() {
				select {
				case logger.ValueChan <- loggerDto.NewLogValue(name, fmt.Sprintf("%v", measure), packet.Timestamp()):
				default:
					break loop
				}
			}

			select {
			case server.PacketChan <- json:
			default:
			}
		}
		
	}()
	server.HandleWebSocketData("/backend/"+os.Getenv("DATA_ENDPOINT"), dataTransferPresentation.DataRoutine)

	go func() {
		for {
			orderChan <- orderTransferInfra.GetPacketValues(<-server.OrderChan)
		}
	}()
	server.HandleWebSocketOrder("/backend/"+os.Getenv("ORDER_ENDPOINT"), orderTransferPresentation.OrderRoutine)

	go func() {
		for {
			server.MessageChan <- messageTransferApplication.NewMessageJSON(messageTransferMappers.GetMessage(<-messageChan))
		}
	}()
	server.HandleWebSocketMessage("/backend/"+os.Getenv("MESSAGE_ENDPOINT"), messageTransferPresentation.MessageRoutine)

	server.HandleLog("/backend/"+os.Getenv("LOG_ENDPOINT"), logger.EnableChan)
	server.HandlePodData("/backend/"+os.Getenv("POD_DATA_ENDPOINT"), serverMappers.NewPodData(boards))
	server.HandlePodData("/backend/"+os.Getenv("ORDER_DESCRIPTION_ENDPOINT"), serverMappers.GetOrders(boards))
	server.HandleSPA()

	packetAdapter.New(config, 10, 10, dataChan, messageChan, orderChan, boards)

	log.Println("Backend Ready!")
	log.Println("\tListening on:", os.Getenv("SERVER_ADDR"))
	go server.ListenAndServe()

	stop := "n"
	for stop == "n" {
		fmt.Scanf("%s", &stop)
	}
}
