package vehicle

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	excel_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/packet_parser"
	"github.com/HyperloopUPV-H8/Backend-H8/pipe"
	"github.com/HyperloopUPV-H8/Backend-H8/sniffer"
	"github.com/HyperloopUPV-H8/Backend-H8/unit_converter"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/internals"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const (
	READ_CHAN_BUF_SIZE  = 100
	ERROR_CHAN_BUF_SIZE = 100
)

type Builder struct {
	sniffer          *sniffer.Sniffer
	config           BuilderConfig
	parser           *packet_parser.PacketParser
	displayConverter *unit_converter.UnitConverter
	podConverter     *unit_converter.UnitConverter
	pipes            map[string]*pipe.Pipe
	idToPipe         map[uint16]string
	trace            zerolog.Logger
}

type BuilderConfig struct {
	TcpClientTag string `toml:"tcp_client_tag"`
	TcpServerTag string `toml:"tcp_server_tag"`
	UdpTag       string `toml:"udp_tag"`
	Network      struct {
		Mtu              uint
		SnifferInterface string `toml:"sniffer_interface"`
	}
}

func NewBuilder(config BuilderConfig) *Builder {
	trace.Trace().Msg("new vehicle builder")
	return &Builder{
		parser:           packet_parser.NewPacketParser(),
		displayConverter: unit_converter.NewUnitConverter("display"),
		podConverter:     unit_converter.NewUnitConverter("pod"),
		sniffer:          nil,
		pipes:            make(map[string]*pipe.Pipe),
		idToPipe:         make(map[uint16]string),
		trace:            trace.With().Str("component", "vehicleBuilder").Logger(),
		config:           config,
	}
}

func (builder *Builder) AddGlobal(global excel_models.GlobalInfo) {
	builder.trace.Debug().Msg("add global")

	var err error
	filter := getFilter(common.Values(global.BoardToIP), global.ProtocolToPort, builder.config.TcpClientTag, builder.config.TcpServerTag, builder.config.UdpTag)
	builder.sniffer, err = sniffer.New(builder.config.Network.SnifferInterface, filter, builder.config.Network)
	if err != nil {
		builder.trace.Fatal().Stack().Err(err).Msg("")
		return
	}

	laddr := common.AddrWithPort(os.Getenv("VEHICLE_LADDR"), global.ProtocolToPort[builder.config.TcpClientTag])
	for board, ip := range global.BoardToIP {
		builder.trace.Debug().Str("board", board).Str("ip", ip).Msg("add board")
		var err error
		builder.pipes[board], err = pipe.New(laddr, common.AddrWithPort(ip, global.ProtocolToPort[builder.config.TcpServerTag]), builder.config.Network.Mtu)
		if err != nil {
			builder.trace.Fatal().Stack().Err(err).Msg("")
			return
		}
	}

	builder.parser.AddGlobal(global)
	builder.displayConverter.AddGlobal(global)
	builder.podConverter.AddGlobal(global)
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

func (builder *Builder) AddPacket(boardName string, packet excel_models.Packet) {
	builder.trace.Debug().Str("id", packet.Description.ID).Str("name", packet.Description.Name).Str("board", boardName).Msg("add packet")
	id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
	if err != nil {
		builder.trace.Error().Stack().Err(err).Msg("")
		return
	}
	builder.idToPipe[uint16(id)] = boardName

	builder.parser.AddPacket(boardName, packet)
	builder.displayConverter.AddPacket(boardName, packet)
	builder.podConverter.AddPacket(boardName, packet)
}

func (builder *Builder) Build() *Vehicle {
	builder.trace.Info().Msg("build")
	vehicle := &Vehicle{
		sniffer:          builder.sniffer,
		parser:           builder.parser,
		displayConverter: builder.displayConverter,
		podConverter:     builder.podConverter,
		pipes:            builder.pipes,

		packetFactory: internals.NewFactory(),

		idToPipe: builder.idToPipe,
		readChan: make(chan []byte, READ_CHAN_BUF_SIZE),
		stats:    newStats(),

		trace: trace.With().Str("component", "vehicle").Logger(),
	}

	vehicle.sniffer.Listen(vehicle.readChan)
	for name, pipe := range vehicle.pipes {
		pipe.SetOutput(vehicle.readChan)
		pipe.OnConnectionChange(func(name string) func(state bool) {
			return func(state bool) {
				vehicle.onConnectionChange(name, state)
			}
		}(name))
	}

	return vehicle
}
