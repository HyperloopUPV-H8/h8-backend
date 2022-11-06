package main

import (
	"fmt"
	"os"
	"time"

	dataTransfer "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer"
	logger "github.com/HyperloopUPV-H8/Backend-H8/Shared/Logger/infra"
	server "github.com/HyperloopUPV-H8/Backend-H8/Shared/Server/infra"
	streaming "github.com/HyperloopUPV-H8/Backend-H8/Shared/StreamingHandlers"
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter"
	excelAdapterDomain "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	excelRetriever "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_retriever"
	transportController "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("./.env")

	ips := []string{"127.0.0.1", "127.0.0.2"}
	document := excelRetriever.GetExcel(os.Getenv("SPREADSHEET_ID"), "excel.xlsx", ".", os.Getenv("SECRET_FILE_PATH"))

	boards := excelAdapter.GetBoards(document)
	packets := make([]excelAdapterDomain.PacketDTO, 0)
	for _, board := range boards {
		packets = append(packets, board.GetPackets()...)
	}

	packetAdapter := transportController.New(ips, packets)

	logger := logger.NewLogger(".", time.Second*5)

	dataTransfer := dataTransfer.New(boards)

	server := server.New(dataTransfer.PacketChannel, make(chan any), make(chan any))

	logger.Run()
	server.HandleLog("/backend/log", logger.EnableChan)

	server.HandleWebSocketData("/backend/data", streaming.DataSocketHandler)

	server.HandleSPA()

	go func() {
		for {
			payload := dataTransfer.Parse(packetAdapter.ReceiveData)
			select {
			case logger.EntryChan <- *payload:
			default:
			}
			dataTransfer.PacketChannel <- payload.Packet
		}
	}()

	fmt.Println("listening")
	server.ListenAndServe()
}
