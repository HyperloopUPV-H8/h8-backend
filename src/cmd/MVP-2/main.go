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
	unitsInfra "github.com/HyperloopUPV-H8/Backend-H8/Shared/units/infra"
	units "github.com/HyperloopUPV-H8/Backend-H8/Shared/units/infra/mappers"
	dataTransferApplication "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/application"
	dataTransfer "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/infra"
	dataTransferPresentation "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/presentation"
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

	packetAdapter := packetAdapter.New(ips, packets)

	podUnitConverter := unitsInfra.NewUnits(podUnits)
	displayUnitConverter := unitsInfra.NewUnits(displayUnits)

	packetFactory := dataTransfer.NewFactory()

	server := createServer()
	logger := logger.NewLogger(".", time.Second*5)

	go func() {
		for {
			packet := units.ConvertUpdate(units.ConvertUpdate(packetAdapter.ReceiveData(), podUnitConverter), displayUnitConverter)
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

	server.HandleLog("/log", logger.EnableChan)
	server.HandleWebSocketData("/data", dataTransferPresentation.DataRoutine)
	server.HandleSPA()

	logger.Run()

	fmt.Println("listening")
	server.ListenAndServe()
}

func createServer() server.HTTPServer[
	dataTransferApplication.PacketJSON,
	orderTransferDomain.OrderWebAdapter,
	messageTransferDomain.Message,
] {
	server := server.New[
		dataTransferApplication.PacketJSON,
		orderTransferDomain.OrderWebAdapter,
		messageTransferDomain.Message]()

	return server
}
