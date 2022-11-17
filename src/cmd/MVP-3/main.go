package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	logger "github.com/HyperloopUPV-H8/Backend-H8/Shared/Logger/infra"
	loggerDto "github.com/HyperloopUPV-H8/Backend-H8/Shared/Logger/infra/dto"
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter"
	excelAdapterDomain "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	excelRetriever "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_retriever"
	packetAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter"
	server "github.com/HyperloopUPV-H8/Backend-H8/Shared/server/infra"
	serverMappers "github.com/HyperloopUPV-H8/Backend-H8/Shared/server/infra/mappers"
	units "github.com/HyperloopUPV-H8/Backend-H8/Shared/units/infra"
	unitsMappers "github.com/HyperloopUPV-H8/Backend-H8/Shared/units/infra/mappers"
	dataTransferApplication "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/application"
	dataTransfer "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/infra"
	dataTransferPresentation "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/infra/presentation"
	messageTransferApplication "github.com/HyperloopUPV-H8/Backend-H8/message_transfer/application"
	messageTransferPresentation "github.com/HyperloopUPV-H8/Backend-H8/message_transfer/infra/presentation"
	orderTransferDomain "github.com/HyperloopUPV-H8/Backend-H8/order_transfer/domain"
	orderTransferInfra "github.com/HyperloopUPV-H8/Backend-H8/order_transfer/infra/mappers"
	orderTransferPresentation "github.com/HyperloopUPV-H8/Backend-H8/order_transfer/infra/presentation"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting Backend...")

	godotenv.Load(".env")

	ips := []string{os.Getenv("TARGET_IP")}

	log.Println("Loading document...")
	document := excelRetriever.GetExcel(os.Getenv("SPREADSHEET_ID"), os.Getenv("SPREADSHEET_NAME"), os.Getenv("SPREADSHEET_PATH"), os.Getenv("CREDENTIALS"))

	boards := excelAdapter.GetBoards(document)
	log.Println("\tDone!")

	log.Println("Fetching packets...")
	packets := getPackets(boards)
	podUnits, displayUnits := getUnits(packets)

	podUnitConverter := units.NewUnits(podUnits)
	displayUnitConverter := units.NewUnits(displayUnits)
	log.Println("\tDone!")

	log.Println("Opening pod connection...")
	packetAdapter := packetAdapter.New(ips, packets)
	log.Println("\tDone!")

	log.Println("Setting up parser...")
	packetFactory := dataTransfer.NewFactory()
	log.Println("\tDone!")

	log.Println("Setting up logger...")
	logger := logger.NewLogger(os.Getenv("LOG_DIR"), time.Second*10)
	log.Println("\tDone!")

	log.Println("Setting up server...")
	server := server.New[
		dataTransferApplication.PacketJSON,
		orderTransferDomain.Order,
		messageTransferApplication.MessageJSON]()

	log.Println("\t\tStarting data routine")
	go func() {
		for {
			packet := unitsMappers.ConvertUpdate(unitsMappers.ConvertUpdate(packetAdapter.ReceiveData(), podUnitConverter), displayUnitConverter)
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
	server.HandleWebSocketData("/backend/packets", dataTransferPresentation.DataRoutine)
	log.Println("\t\t\tDone!")

	log.Println("\t\tStarting order routine")
	go func() {
		for {
			rawOrder := <-server.OrderChan
			order := orderTransferInfra.GetPacketValues(rawOrder)
			packetAdapter.Send(os.Getenv("TARGET_IP"), order)
		}
	}()
	server.HandleWebSocketOrder("/backend/order", orderTransferPresentation.OrderRoutine)
	log.Println("\t\t\tDone!")

	log.Println("\t\tStarting message routine")
	go func() {

	}()
	server.HandleWebSocketMessage("/backend/message", messageTransferPresentation.MessageRoutine)
	log.Println("\t\t\tDone!")

	server.HandleLog("/backend/log", logger.EnableChan)
	server.HandlePodData("/backend/podDataDescription", serverMappers.NewPodData(boards))
	server.HandleSPA()

	log.Println("\tDone!")

	log.Println("Backend Ready!")
	log.Println("\tListening on:", os.Getenv("SERVER_ADDR"))
	server.ListenAndServe()
}

func getPackets(boards map[string]excelAdapterDomain.BoardDTO) []excelAdapterDomain.PacketDTO {
	packets := make([]excelAdapterDomain.PacketDTO, 0)
	for _, board := range boards {
		packets = append(packets, board.GetPackets()...)
	}
	return packets
}

func getUnits(packets []excelAdapterDomain.PacketDTO) (pod map[string]string, display map[string]string) {
	pod = make(map[string]string)
	display = make(map[string]string)
	for _, packet := range packets {
		for _, measurement := range packet.Measurements {
			pod[measurement.Name] = parseUnits(measurement.PodUnits)
			display[measurement.Name] = parseUnits(measurement.DisplayUnits)
		}
	}
	return pod, display
}

func parseUnits(units string) string {
	if split := strings.Split(units, "#"); len(split) > 0 {
		return split[1]
	} else {
		return ""
	}
}
