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
	packetParserInfra "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra"
	transportControllerInfra "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/transport_controller/infra"
	unitsInfra "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/units/infra"
	unitsMappers "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/units/infra/mappers"
	server "github.com/HyperloopUPV-H8/Backend-H8/Shared/server/infra"
	serverMappers "github.com/HyperloopUPV-H8/Backend-H8/Shared/server/infra/mappers"
	dataTransferApplication "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/application"
	dataTransfer "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/infra"
	dataTransferPresentation "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/infra/presentation"
	messageTransferApplication "github.com/HyperloopUPV-H8/Backend-H8/message_transfer/application"
	orderTransferApplication "github.com/HyperloopUPV-H8/Backend-H8/order_transfer/application"
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

	transportController := transportControllerInfra.NewTransportController(config)
	packetParser := packetParserInfra.NewPacketAggregate(boards)
	podUnits := unitsInfra.NewPodUnitAggregate(boards)
	displayUnits := unitsInfra.NewDisplayUnitAggregate(boards)

	packetFactory := dataTransfer.NewFactory()

	logger := logger.NewLogger(os.Getenv("LOG_DIR"), time.Second*10)

	server := server.New[
		dataTransferApplication.PacketJSON,
		orderTransferApplication.OrderJSON,
		messageTransferApplication.MessageJSON](1024)

	count1 := 0
	count2 := 0
	go func() {
		for {
			packet, _ := transportController.ReceiveData()
			count1 += 1
			dto := packetParserInfra.Decode(packet, *packetParser)
			unitsMappers.ConvertUpdate(&dto, *podUnits)
			unitsMappers.ConvertUpdate(&dto, *displayUnits)

			json := dataTransferApplication.NewJSON(packetFactory.NewPacket(dto))

		loop:
			for name, measure := range dto.Values() {
				select {
				case logger.ValueChan <- loggerDto.NewLogValue(name, fmt.Sprintf("%v", measure), dto.Timestamp()):
				default:
					break loop
				}
			}

			select {
			case server.PacketChan <- json:
				count2 += 1
			default:
			}
		}

	}()
	server.HandleWebSocketData("/backend/"+os.Getenv("DATA_ENDPOINT"), dataTransferPresentation.DataRoutine)

	server.HandleLog("/backend/"+os.Getenv("LOG_ENDPOINT"), logger.EnableChan)
	server.HandlePodData("/backend/"+os.Getenv("POD_DATA_ENDPOINT"), serverMappers.NewPodData(boards))
	server.HandlePodData("/backend/"+os.Getenv("ORDER_DESCRIPTION_ENDPOINT"), serverMappers.GetOrders(boards))
	server.HandleSPA()

	log.Println("Backend Ready!")
	log.Println("\tListening on:", os.Getenv("SERVER_ADDR"))
	go server.ListenAndServe()

	stop := "n"
	for stop == "n" {
		fmt.Scanf("%s", &stop)
	}
	log.Println(count2, "/", count1)
}
