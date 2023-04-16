package vehicle

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	excel_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/message_parser"
	"github.com/HyperloopUPV-H8/Backend-H8/packet_parser"
	"github.com/HyperloopUPV-H8/Backend-H8/pipe"
	"github.com/HyperloopUPV-H8/Backend-H8/sniffer"
	"github.com/HyperloopUPV-H8/Backend-H8/unit_converter"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/internals"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const (
	UPDATE_CHAN_BUF_SIZE  = 100
	MESSAGE_CHAN_BUF_SIZE = 100
	ERROR_CHAN_BUF_SIZE   = 100
)

type VehicleConfig struct {
	TcpClientTag string `toml:"tcp_client_tag"`
	TcpServerTag string `toml:"tcp_server_tag"`
	UdpTag       string `toml:"udp_tag"`
	Network      struct {
		Mtu              uint
		SnifferInterface string `toml:"sniffer_interface"`
	}
	Messages message_parser.MessageParserConfig
}

func NewVehicle(boards map[string]excel_models.Board, globalInfo excel_models.GlobalInfo, config VehicleConfig, onConnectionChange func(string, bool)) Vehicle {
	trace.Trace().Msg("new vehicle builder")

	messageChan := make(chan []byte, MESSAGE_CHAN_BUF_SIZE)
	trace := trace.With().Str("component", "vehicle").Logger()
	vehicle := Vehicle{
		parser:             packet_parser.New(boards),
		messageParser:      message_parser.New(globalInfo, config.Messages),
		displayConverter:   unit_converter.NewUnitConverter("display", boards, globalInfo.UnitToOperations),
		podConverter:       unit_converter.NewUnitConverter("pod", boards, globalInfo.UnitToOperations),
		sniffer:            createSniffer(globalInfo, config, trace),
		pipes:              createPipes(globalInfo, messageChan, onConnectionChange, config, trace),
		idToBoard:          getIdToBoard(boards, trace),
		packetFactory:      internals.NewFactory(),
		updateChan:         make(chan []byte, UPDATE_CHAN_BUF_SIZE),
		messageChan:        messageChan,
		onConnectionChange: onConnectionChange,
		stats:              newStats(),
		trace:              trace,
	}

	vehicle.sniffer.Listen(vehicle.updateChan)

	return vehicle
}

func createSniffer(global excel_models.GlobalInfo, config VehicleConfig, trace zerolog.Logger) sniffer.Sniffer {
	filter := getFilter(common.Values(global.BoardToIP), global.ProtocolToPort, config.TcpClientTag, config.TcpServerTag, config.UdpTag)
	sniffer, err := sniffer.New(filter, config.Network)

	if err != nil {
		trace.Fatal().Stack().Err(err).Msg("error creating sniffer")
	}

	return *sniffer
}

func createPipes(global excel_models.GlobalInfo, messageChan chan []byte, onConnectionChange func(string, bool), config VehicleConfig, trace zerolog.Logger) map[string]*pipe.Pipe {
	laddr := common.AddrWithPort(global.BackendIP, global.ProtocolToPort[config.TcpClientTag])
	pipes := make(map[string]*pipe.Pipe)
	for board, ip := range global.BoardToIP {
		raddr := common.AddrWithPort(ip, global.ProtocolToPort[config.TcpServerTag])
		pipe, err := pipe.New(laddr, raddr, config.Network.Mtu, messageChan, func(state bool) { onConnectionChange(board, state) })
		if err != nil {
			trace.Fatal().Stack().Err(err).Msg("error creating pipe")

		}

		pipes[board] = pipe
	}
	return pipes
}

func getFilter(addrs []string, protocolToPort map[string]string, tcpClientTag string, tcpServerTag string, udpTag string) string {
	udp := fmt.Sprintf("(udp port %s)", protocolToPort[udpTag])
	udpAddr := ""
	for _, addr := range addrs {
		udpAddr = fmt.Sprintf("%s or (src host %s)", udpAddr, addr)
	}
	udpAddr = strings.TrimPrefix(udpAddr, " or ")
	udp = fmt.Sprintf("%s and (%s)", udp, udpAddr)

	tcp := fmt.Sprintf("(tcp port %s or tcp port %s) and (tcp[tcpflags] & (tcp-fin | tcp-syn | tcp-ack) == 0)", protocolToPort[tcpClientTag], protocolToPort[tcpServerTag])
	tcpAddrSrc := ""
	tcpAddrDst := ""
	for _, addr := range addrs {
		tcpAddrSrc = fmt.Sprintf("%s or (src host %s)", tcpAddrSrc, addr)
		tcpAddrDst = fmt.Sprintf("%s or (dst host %s)", tcpAddrDst, addr)
	}
	tcpAddrSrc = strings.TrimPrefix(tcpAddrSrc, " or ")
	tcpAddrDst = strings.TrimPrefix(tcpAddrDst, " or ")
	tcp = fmt.Sprintf("%s and (%s) and (%s)", tcp, tcpAddrSrc, tcpAddrDst)

	filter := fmt.Sprintf("(%s) or (%s)", udp, tcp)
	trace.Trace().Strs("addrs", addrs).Str("filter", filter).Msg("new filter")
	return filter
}

func getIdToBoard(boards map[string]excel_models.Board, trace zerolog.Logger) map[uint16]string {
	idToBoard := make(map[uint16]string)
	for _, board := range boards {
		for _, packet := range board.Packets {
			id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
			if err != nil {
				trace.Fatal().Stack().Err(err).Msg("error parsing id")
				continue
			}
			idToBoard[uint16(id)] = board.Name
		}
	}

	return idToBoard
}
