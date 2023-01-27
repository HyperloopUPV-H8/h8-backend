package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/connection_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer"
	dataTransferModels "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/models"
	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter"
	"github.com/HyperloopUPV-H8/Backend-H8/log_handle"
	logHandleModels "github.com/HyperloopUPV-H8/Backend-H8/log_handle/models"
	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/order_transfer"
	orderTransferModels "github.com/HyperloopUPV-H8/Backend-H8/order_transfer/models"
	"github.com/HyperloopUPV-H8/Backend-H8/packet_parser"
	"github.com/HyperloopUPV-H8/Backend-H8/server"
	"github.com/HyperloopUPV-H8/Backend-H8/transport_controller"
	transportControllerModels "github.com/HyperloopUPV-H8/Backend-H8/transport_controller/models"
	"github.com/HyperloopUPV-H8/Backend-H8/unit_converter"
	"github.com/HyperloopUPV-H8/Backend-H8/websocket_handle"
	"github.com/HyperloopUPV-H8/Backend-H8/websocket_handle/models"
	"github.com/google/gopacket/pcap"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	godotenv.Load(".env")

	document := excel_adapter.FetchDocument(os.Getenv("EXCEL_ID"), os.Getenv("EXCEL_PATH"), os.Getenv("EXCEL_NAME"))
	controlSections := excel_adapter.GetControlSections(document)
	controlSectionsRaw, err := json.Marshal(controlSections)
	if err != nil {
		log.Fatalln(err)
	}

	podConverter := unit_converter.UnitConverter{Kind: "pod"}
	displayConverter := unit_converter.UnitConverter{Kind: "display"}

	packetParser := packet_parser.NewPacketParser()

	podData := dataTransferModels.PodData{Boards: make(map[string]dataTransferModels.Board)}
	podDataRaw, err := json.Marshal(podData)
	if err != nil {
		log.Fatalln(err)
	}

	orderData := orderTransferModels.OrderData{}
	orderDataRaw, err := json.Marshal(orderData)
	if err != nil {
		log.Fatalln(err)
	}

	idToType := IDtoType{}
	idToIP := IDtoIP{}
	ipToBoard := IPtoBoard{}

	excel_adapter.AddExpandedPackets(document, &podConverter, &displayConverter, &packetParser, &podData, &orderData, &idToType, &idToIP, &ipToBoard)

	laddr, err := net.ResolveTCPAddr("tcp", os.Getenv("LOCAL_ADDRESS"))
	if err != nil {
		log.Fatalln(err)
	}

	rawRAddrs := strings.Split(os.Getenv("REMOTE_ADDRESSES"), ",")
	raddrs := make([]*net.TCPAddr, len(rawRAddrs))
	for i, rawRAddr := range rawRAddrs {
		raddr, err := net.ResolveTCPAddr("tcp", rawRAddr)
		if err != nil {
			log.Fatalln(err)
		}
		raddrs[i] = raddr
	}

	connectionTransfer, connectionChannel := connection_transfer.New()

	dataTransfer, dataTransferChannel := data_transfer.New(time.Millisecond * 10)

	messageTransfer, messageChannel := message_transfer.New()

	orderChannel := make(chan orderTransferModels.Order, 100)
	orderTransfer, ordChannel := order_transfer.New(orderChannel)

	packetFactory := data_transfer.NewFactory()

	transportControllerConfig := transportControllerModels.Config{
		Dump:    make(chan []byte),
		Snaplen: 2000,
		Promisc: true,
		Timeout: pcap.BlockForever,
		BPF:     getFilter(raddrs),
		OnConnUpdate: func(addr *net.TCPAddr, up bool) {
			connectionTransfer.Update(ipToBoard[addr.IP.String()], up)
		},
	}

	live, err := strconv.ParseBool(os.Getenv("SNIFFER_LIVE"))
	if err != nil {
		log.Fatalln(err)
	}

	transportController := transport_controller.Open(laddr, raddrs, os.Getenv("SNIFFER_DEV"), live, transportControllerConfig)
	defer transportController.Close()

	logger, loggerChannel := log_handle.NewLogger(logHandleModels.Config{
		DumpSize: 7000,
		RowSize:  20,
		BasePath: os.Getenv("LOG_PATH"),
		Updates:  make(chan map[string]any, 10000),
		Autosave: time.NewTicker(time.Minute),
	})

	go func(msgChannel chan models.MessageTarget) {
		for packet := range transportControllerConfig.Dump {
			id, values := packetParser.Decode(packet)
			values = podConverter.Convert(values)
			values = displayConverter.Convert(values)
			logger.Update(values)
			if idToType[id] == "data" {
				dataTransfer.Update(packetFactory.NewPacket(id, packet, values))
			} else {
				messageTransfer.Broadcast(packetFactory.NewPacket(id, packet, values))
			}
		}
	}(messageChannel)

	go func() {
		for order := range orderChannel {
			order.Values = displayConverter.Revert(order.Values)
			order.Values = podConverter.Revert(order.Values)
			transportController.Write(idToIP[order.ID], packetParser.Encode(order.ID, order.Values))
		}
	}()

	httpServer := server.Server{Router: mux.NewRouter()}
	httpServer.ServeData("/backend/"+os.Getenv("POD_DATA_ENDPOINT"), podDataRaw)
	httpServer.ServeData("/backend/"+os.Getenv("CONTROL_SECTIONS_ENDPOINT"), controlSectionsRaw)
	httpServer.ServeData("/backend/"+os.Getenv("ORDER_DATA_ENDPOINT"), orderDataRaw)

	handle := websocket_handle.RunWSHandle(httpServer.Router, "/backend", map[string]chan models.MessageTarget{
		"podData/update":    dataTransferChannel,
		"message/update":    messageChannel,
		"order/update":      ordChannel,
		"connection/update": connectionChannel,
		"logger/enable":     loggerChannel,
	})

	path, _ := os.Getwd()
	httpServer.FileServer("/", filepath.Join(path, "static"))

	go httpServer.ListenAndServe(os.Getenv("SERVER_ADDR"))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

loop:
	for {
		select {
		case <-time.After(time.Second * 10):
			fmt.Println(transportController.Stats())
		case <-interrupt:
			break loop
		}
	}

	log.Println(handle, orderTransfer)
}
