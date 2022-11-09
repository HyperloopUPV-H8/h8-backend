package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	loggerDto "github.com/HyperloopUPV-H8/Backend-H8/Shared/Logger/infra/dto"
	server "github.com/HyperloopUPV-H8/Backend-H8/Shared/Server/infra"
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter"
	excelAdapterDomain "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	excelRetriever "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_retriever"
	logger "github.com/HyperloopUPV-H8/Backend-H8/Shared/logger/infra"
	packetAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter"
	streaming "github.com/HyperloopUPV-H8/Backend-H8/Shared/streaming_handlers"
	unitsInfra "github.com/HyperloopUPV-H8/Backend-H8/Shared/units/infra"
	units "github.com/HyperloopUPV-H8/Backend-H8/Shared/units/infra/mappers"
	dataTransfer "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/infra"
	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/infra/dto"
	messageTransferDomain "github.com/HyperloopUPV-H8/Backend-H8/message_transfer/domain"
	orderTransferDomain "github.com/HyperloopUPV-H8/Backend-H8/order_transfer/domain"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting program")
	godotenv.Load(".env")

	ips := []string{"127.0.0.1", "127.0.0.2"}
	document := excelRetriever.GetExcel(os.Getenv("SPREADSHEET_ID"), "excel.xlsx", ".", os.Getenv("SECRET_FILE_PATH"))

	boards := excelAdapter.GetBoards(document)
	packets := make([]excelAdapterDomain.PacketDTO, 0)
	podUnits := make(map[string]string)
	displayUnits := make(map[string]string)
	for _, board := range boards {
		expandedPackets := board.GetPackets()
		packets = append(packets, expandedPackets...)
		for _, packet := range expandedPackets {
			for _, value := range packet.Measurements {
				podUnits[value.Name] = strings.Split(value.PodUnits, "#")[1]
				displayUnits[value.Name] = strings.Split(value.DisplayUnits, "#")[1]
			}

		}
	}

	// TODO: Change main logic to incorporate structure changes
	packetAdapter := packetAdapter.New(ips, packets)

	unitConverter := unitsInfra.NewUnits(podUnits, displayUnits)

	logger := logger.NewLogger(".", time.Second*5)
	server := createServer()

	server.HandleLog("/log", logger.EnableChan)
	server.HandleWebSocketData("/data", streaming.DataSocketHandler)
	server.HandleSPA()

	packetFactory := dataTransfer.NewFactory()

	go func() {
		raw := packetAdapter.ReceiveData()
		update := units.ConvertUpdate(raw, unitConverter)
		fmt.Println(raw)
		fmt.Println(update)
		for {
			raw := packetAdapter.ReceiveData()
			update := units.ConvertUpdate(raw, unitConverter)
			packet := packetFactory.NewPacket(update)

			select {
			case server.PacketChan <- packet:
			default:
			}

		loop:
			for name, measure := range packet.Values() {
				select {
				case logger.ValueChan <- loggerDto.NewLogValue(name, fmt.Sprintf("%v", measure), update.Timestamp()):
				default:
					break loop
				}
			}
		}
	}()

	logger.Run()

	fmt.Println("listening")
	server.ListenAndServe()
}

func createServer() server.HTTPServer[
	dto.Packet,
	orderTransferDomain.OrderWebAdapter,
	messageTransferDomain.Message,
] {
	server := server.New[
		dto.Packet,
		orderTransferDomain.OrderWebAdapter,
		messageTransferDomain.Message]()

	return server
}
