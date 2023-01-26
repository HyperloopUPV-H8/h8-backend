package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer"
	dataTransferModels "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/models"
	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter"
	"github.com/HyperloopUPV-H8/Backend-H8/log_handle"
	logHandleModels "github.com/HyperloopUPV-H8/Backend-H8/log_handle/models"
	orderTransferModels "github.com/HyperloopUPV-H8/Backend-H8/order_transfer/models"
	"github.com/HyperloopUPV-H8/Backend-H8/packet_parser"
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

	podConverter := unit_converter.UnitConverter{Kind: "pod"}
	displayConverter := unit_converter.UnitConverter{Kind: "display"}

	packetParser := packet_parser.NewPacketParser()

	podData := dataTransferModels.PodData{Boards: make(map[string]dataTransferModels.Board)}
	orderData := orderTransferModels.OrderData{}

	idToType := IDtoType{}
	idToIP := IDtoIP{}
	ipToBoard := IPtoBoard{}

	excel_adapter.Compile(document, &podConverter, &displayConverter, &packetParser, &podData, &orderData, &idToType, &idToIP, &ipToBoard)

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

	dataTransfer := data_transfer.New(time.Millisecond * 10)
	defer dataTransfer.Close()

	packetFactory := data_transfer.NewFactory()

	transportControllerConfig := transportControllerModels.Config{
		Dump:    make(chan []byte),
		Snaplen: 2000,
		Promisc: true,
		Timeout: pcap.BlockForever,
		BPF:     getFilter(raddrs),
		OnConnUpdate: func(addr *net.TCPAddr, up bool) {
		},
	}

	live, err := strconv.ParseBool(os.Getenv("SNIFFER_LIVE"))
	if err != nil {
		log.Fatalln(err)
	}

	transportController := transport_controller.Open(laddr, raddrs, os.Getenv("SNIFFER_DEV"), live, transportControllerConfig)
	defer transportController.Close()

	logger := log_handle.NewLogger(logHandleModels.Config{
		DumpSize: 7000,
		RowSize:  20,
		BasePath: os.Getenv("LOG_PATH"),
		Updates:  make(chan map[string]any, 10000),
		Autosave: time.NewTicker(time.Minute),
	})

	router := mux.NewRouter()

	packetUpdateChannel := make(chan models.MessageTarget)
	podDataChannel := make(chan models.MessageTarget)
	handler := websocket_handle.RunWSHandle(router, "backend", map[string]chan models.MessageTarget{
		"podData":      podDataChannel,
		"packetUpdate": packetUpdateChannel,
	})

	go func(channel chan models.MessageTarget) {
		for packet := range transportControllerConfig.Dump {
			id, values := packetParser.Decode(packet)
			values = podConverter.Convert(values)
			values = displayConverter.Convert(values)
			logger.Update(values)
			if idToType[id] == "data" {
				packet := packetFactory.NewPacket(id, packet, values)
				channel <- models.MessageTarget{
					Target: handler.GetClients(),
					Msg: models.Message{
						Kind: "podData/update",
						Msg:  []dataTransferModels.PacketUpdate{packet},
					},
				}
			}
		}
	}(packetUpdateChannel)

	go func(channel chan models.MessageTarget) {
		for msg := range channel {
			channel <- models.MessageTarget{
				Target: msg.Target,
				Msg: models.Message{
					Kind: "podData/structure",
					Msg:  podData,
				},
			}
		}
	}(podDataChannel)

	router.PathPrefix("/").HandlerFunc(http.FileServer(http.Dir(path.Join(".", "static"))).ServeHTTP)

	go log.Fatalln(http.ListenAndServe("127.0.0.1:4000", router))

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
	log.Println(handler)
}
