package main

import (
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

	ips := []string{}
	document := excelRetriever.GetExcel("excel.xlsx", ".")

	boards := excelAdapter.GetBoards(document)
	packets := make([]excelAdapterDomain.PacketDTO, 0)
	for _, board := range boards {
		packets = append(packets, board.GetPackets()...)
	}

	packetAdapter := transportController.New(ips, packets)

	logger := logger.New(".", time.Second*5)

	dataTransfer := dataTransfer.New(boards)

	server := server.New(dataTransfer.PacketChannel, make(chan any), make(chan any))

	server.HandleLog("/backend/log", logger.EnableChan)
	logger.Run()

	dataTransfer.Invoke(packetAdapter.ReceiveData)

	server.HandleWebSocketData("/backend/data", streaming.DataSocketHandler)
}
