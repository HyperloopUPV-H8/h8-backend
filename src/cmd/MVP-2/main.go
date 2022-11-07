package main

import (
	"fmt"
	"os"
	"time"

	server "github.com/HyperloopUPV-H8/Backend-H8/Shared/Server/infra"
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter"
	excelAdapterDomain "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	excelRetriever "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_retriever"
	logger "github.com/HyperloopUPV-H8/Backend-H8/Shared/logger/infra"
	packetAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter"
	streaming "github.com/HyperloopUPV-H8/Backend-H8/Shared/streaming_handlers"
	dataTransfer "github.com/HyperloopUPV-H8/Backend-H8/data_transfer"
	dataTransferDomain "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/domain"
	messageTransferDomain "github.com/HyperloopUPV-H8/Backend-H8/message_transfer/domain"
	orderTransferDomain "github.com/HyperloopUPV-H8/Backend-H8/order_transfer/domain"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting program")
	godotenv.Load("../../.env")

	ips := []string{"127.0.0.1", "127.0.0.2"}
	document := excelRetriever.GetExcel(os.Getenv("SPREADSHEET_ID"), "excel.xlsx", ".", os.Getenv("SECRET_FILE_PATH"))

	boards := excelAdapter.GetBoards(document)
	packets := make([]excelAdapterDomain.PacketDTO, 0)
	for _, board := range boards {
		packets = append(packets, board.GetPackets()...)
	}

	packetAdapter := packetAdapter.New(ips, packets)

	dataTransfer := dataTransfer.New(boards)
	dataTransfer.Invoke(packetAdapter.ReceiveData)

	logger := logger.NewLogger(".", time.Second*5)
	server := createServer()

	server.HandleLog("/log", logger.EnableChan)
	server.HandleWebSocketData("/data", streaming.DataSocketHandler)
	server.HandleWebSocketOrder("/order", streaming.OrderSocketHandler)
	//server.HandleWebSocketMessage("/message", streaming.MessageSocketHandler)
	server.HandleSPA()

	go func() {
		for packetTimestampPair := range dataTransfer.PacketTimestampChannel {
			select {
			case logger.EntryChan <- packetTimestampPair:
			default:
			}
			select {
			case server.PacketChan <- packetTimestampPair.Packet:
			default:
			}

		}
	}()

	//logger.Run()

	// go func() {
	// 	for {
	// 		payload := dataTransfer.Parse(packetAdapter.ReceiveData)
	// 		select {
	// 		case logger.EntryChan <- *payload:
	// 		default:
	// 		}
	// 		dataTransfer.PacketTimestampChan <- payload.Packet
	// 	}
	// }()

	fmt.Println("listening")
	server.ListenAndServe()
}

func createServer() server.HTTPServer[
	dataTransferDomain.Packet,
	orderTransferDomain.OrderWebAdapter,
	messageTransferDomain.Message,
] {
	server := server.New[
		dataTransferDomain.Packet,
		orderTransferDomain.OrderWebAdapter,
		messageTransferDomain.Message]()

	return server
}
