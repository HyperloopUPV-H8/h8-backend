package vehicle

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	excel_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/HyperloopUPV-H8/Backend-H8/pipe"
	"github.com/HyperloopUPV-H8/Backend-H8/sniffer"
	"github.com/HyperloopUPV-H8/Backend-H8/unit_converter"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/packet_parser"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const (
	UPDATE_CHAN_BUF_SIZE = 100
)

type VehicleConstructorArgs struct {
	Boards             map[string]excel_models.Board
	GlobalInfo         excel_models.GlobalInfo
	Config             Config
	OnConnectionChange func(string, bool)
}

func New(args VehicleConstructorArgs) Vehicle {
	trace.Trace().Msg("creating vehicle")

	vehicleTrace := trace.With().Str("component", "vehicle").Logger()
	dataChan := make(chan packet.Packet, UPDATE_CHAN_BUF_SIZE)

	packetParser, err := createPacketParser(args.GlobalInfo, args.Boards, vehicleTrace)

	if err != nil {
		vehicleTrace.Fatal().Err(err).Msg("error creating packetParser")
	}

	names, err := getPacketToValuesNames(args.GlobalInfo, args.Boards)

	if err != nil {
		vehicleTrace.Error().Err(err).Msg("error getting packet to values names")
	}

	vehicle := Vehicle{
		podConverter:     unit_converter.NewUnitConverter("pod", args.Boards, args.GlobalInfo.UnitToOperations),
		displayConverter: unit_converter.NewUnitConverter("display", args.Boards, args.GlobalInfo.UnitToOperations),

		sniffer: createSniffer(args.GlobalInfo, args.Config, vehicleTrace),
		pipes:   createPipes(args.GlobalInfo, dataChan, args.OnConnectionChange, args.Config, vehicleTrace),

		packetParser:     packetParser,
		protectionParser: NewProtectionParser(args.GlobalInfo, args.Config.Protections),
		bitarrayParser:   NewBitarrayParser(names),

		dataChan: dataChan,

		idToBoard:          getIdToBoard(args.Boards, vehicleTrace),
		onConnectionChange: args.OnConnectionChange,
		trace:              vehicleTrace,
	}

	vehicle.sniffer.Listen(dataChan)

	return vehicle
}

func createSniffer(global excel_models.GlobalInfo, config Config, trace zerolog.Logger) sniffer.Sniffer {
	filter := getFilter(common.Values(global.BoardToIP), global.ProtocolToPort, config.Network.TcpClientTag, config.Network.TcpServerTag, config.Network.UdpTag)
	sniffer, err := sniffer.New(filter, config.Network.Sniffer)

	if err != nil {
		trace.Fatal().Stack().Err(err).Msg("error creating sniffer")
	}

	return *sniffer
}

func createPipes(global excel_models.GlobalInfo, dataChan chan<- packet.Packet, onConnectionChange func(string, bool), config Config, trace zerolog.Logger) map[string]*pipe.Pipe {
	laddr := common.AddrWithPort(global.BackendIP, global.ProtocolToPort[config.Network.TcpClientTag])
	pipes := make(map[string]*pipe.Pipe)
	for board, ip := range global.BoardToIP {
		raddr := common.AddrWithPort(ip, global.ProtocolToPort[config.Network.TcpServerTag])
		pipe, err := pipe.New(laddr, raddr, config.Network.Sniffer.Mtu, dataChan, getOnConnectionChange(board, onConnectionChange))
		if err != nil {
			trace.Fatal().Stack().Err(err).Msg("error creating pipe")
		}

		pipes[board] = pipe
	}
	return pipes
}

func getOnConnectionChange(board string, onConnectionChange func(string, bool)) func(bool) {
	return func(state bool) {
		onConnectionChange(board, state)
	}
}

func createPacketParser(global excel_models.GlobalInfo, boards map[string]excel_models.Board, trace zerolog.Logger) (packet_parser.PacketParser, error) {
	structures, err := getStructures(global, boards)
	if err != nil {
		return packet_parser.PacketParser{}, err
	}

	ids := getDataIds(global, boards, trace)

	return packet_parser.NewPacketParser(ids, structures, getEnums(global, boards)), nil
}

func getStructures(global excel_models.GlobalInfo, boards map[string]excel_models.Board) (map[uint16][]packet.ValueDescriptor, error) {
	structures := make(map[uint16][]packet.ValueDescriptor)
	for _, board := range boards {
		for _, packet := range board.Packets {
			if packet.Description.Type == "data" || packet.Description.Type == "order" {
				id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
				if err != nil {
					return nil, err
				}
				structures[uint16(id)] = getDescriptor(packet.Values)
			}
		}
	}
	return structures, nil
}

func getDataIds(global excel_models.GlobalInfo, boards map[string]excel_models.Board, trace zerolog.Logger) common.Set[uint16] {
	ids := common.NewSet[uint16]()
	for _, board := range boards {
		for _, packet := range board.Packets {
			if packet.Description.Type == "data" {
				id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
				if err != nil {
					trace.Error().Err(err).Msg("error parsing packet id")
					continue
				}
				ids.Add(uint16(id))
			}
		}
	}
	return ids
}

func getDescriptor(values []excel_models.Value) []packet.ValueDescriptor {
	descriptor := make([]packet.ValueDescriptor, len(values))
	for i, value := range values {
		descriptor[i] = packet.ValueDescriptor{
			Name: value.ID,
			Type: getValueType(value.Type),
		}
	}
	return descriptor
}

func getValueType(literal string) string {
	if strings.HasPrefix(literal, "enum") {
		return "enum"
	} else {
		return literal
	}
}

func getEnums(global excel_models.GlobalInfo, boards map[string]excel_models.Board) map[string]packet.EnumDescriptor {
	enums := make(map[string]packet.EnumDescriptor)
	for _, board := range boards {
		for _, packet := range board.Packets {
			for _, value := range packet.Values {
				if getValueType(value.Type) != "enum" {
					continue
				}
				enums[value.ID] = getEnumDescriptor(value.Type)
			}
		}
	}
	return enums
}

func getEnumDescriptor(literal string) packet.EnumDescriptor {
	withoutSpaceLiteral := strings.ReplaceAll(literal, " ", "")
	optionsLiteral := strings.TrimSuffix(strings.TrimPrefix(withoutSpaceLiteral, "enum("), ")")
	return strings.Split(optionsLiteral, ",")
}

func getPacketToValuesNames(global excel_models.GlobalInfo, boards map[string]excel_models.Board) (map[uint16][]string, error) {
	names := make(map[uint16][]string)
	for _, board := range boards {
		for _, packet := range board.Packets {
			id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
			if err != nil {
				return nil, err
			}
			names[(uint16)(id)] = getNamesFromValues(packet.Values)
		}
	}

	return names, nil
}

func getNamesFromValues(values []excel_models.Value) []string {
	names := make([]string, len(values))
	for i, value := range values {
		names[i] = value.ID
	}
	return names
}

func getFilter(addrs []string, protocolToPort map[string]string, tcpClientTag string, tcpServerTag string, udpTag string) string {
	ipip := "ip[9] == 4"

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

	filter := fmt.Sprintf("(%s) or (%s) or (%s)", ipip, udp, tcp)
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
