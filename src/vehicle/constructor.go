package vehicle

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	excel_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/data"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/message"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/order"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/parsers"
	"github.com/HyperloopUPV-H8/Backend-H8/pipe"
	"github.com/HyperloopUPV-H8/Backend-H8/sniffer"
	"github.com/HyperloopUPV-H8/Backend-H8/unit_converter"
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
	trace.Trace().Msg("new vehicle builder")

	vehicleTrace := trace.With().Str("component", "vehicle").Logger()
	dataChan := make(chan packet.Raw, UPDATE_CHAN_BUF_SIZE)
	parser, err := createParser(args.GlobalInfo, args.Boards)
	if err != nil {
		// FIXME: handle error
		panic(err)
	}
	vehicle := Vehicle{
		parser:             parser,
		podConverter:       unit_converter.NewUnitConverter("pod", args.Boards, args.GlobalInfo.UnitToOperations),
		displayConverter:   unit_converter.NewUnitConverter("display", args.Boards, args.GlobalInfo.UnitToOperations),
		sniffer:            createSniffer(args.GlobalInfo, args.Config, vehicleTrace),
		pipes:              createPipes(args.GlobalInfo, dataChan, args.OnConnectionChange, args.Config, vehicleTrace),
		idToBoard:          getIdToBoard(args.Boards, vehicleTrace),
		dataChan:           dataChan,
		onConnectionChange: args.OnConnectionChange,
		trace:              vehicleTrace,
	}

	vehicle.sniffer.Listen(dataChan)

	return vehicle
}

func createSniffer(global excel_models.GlobalInfo, config Config, trace zerolog.Logger) sniffer.Sniffer {
	filter := getFilter(common.Values(global.BoardToIP), global.ProtocolToPort, config.TcpClientTag, config.TcpServerTag, config.UdpTag)
	sniffer, err := sniffer.New(filter, config.Network)

	if err != nil {
		trace.Fatal().Stack().Err(err).Msg("error creating sniffer")
	}

	return *sniffer
}

func createPipes(global excel_models.GlobalInfo, dataChan chan<- packet.Raw, onConnectionChange func(string, bool), config Config, trace zerolog.Logger) map[string]*pipe.Pipe {
	laddr := common.AddrWithPort(global.BackendIP, global.ProtocolToPort[config.TcpClientTag])
	pipes := make(map[string]*pipe.Pipe)
	for board, ip := range global.BoardToIP {
		raddr := common.AddrWithPort(ip, global.ProtocolToPort[config.TcpServerTag])
		// FIXME: func(state bool) does not work (closure takes the same board)
		pipe, err := pipe.New(laddr, raddr, config.Network.Mtu, dataChan, getOnConnectionChange(board, onConnectionChange))
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

func createParser(global excel_models.GlobalInfo, boards map[string]excel_models.Board) (*packet.Parser, error) {
	structures, err := getStructures(global, boards)
	if err != nil {
		return nil, err
	}

	valueParser := parsers.NewValueParser(structures, getEnums(global, boards))

	names, err := getNames(global, boards)
	if err != nil {
		return nil, err
	}
	bitarrayParser := parsers.NewBitarrayParser(names)

	dataParser := data.NewParser(valueParser)

	config, err := getMessageConfig(global, boards)
	if err != nil {
		return nil, err
	}

	messageParser := message.NewParser(config)
	orderParser := order.NewParser(valueParser, bitarrayParser)

	kinds, err := getIdKinds(global, boards)
	if err != nil {
		return nil, err
	}

	return packet.NewParser(
		kinds,
		map[packet.Kind]packet.Decoder{
			packet.Data:    dataParser,
			packet.Message: messageParser,
			packet.Order:   orderParser,
		},
		map[packet.Kind]packet.Encoder{
			packet.Order: orderParser,
		},
	), nil
}

func getStructures(global excel_models.GlobalInfo, boards map[string]excel_models.Board) (map[uint16][]packet.ValueDescriptor, error) {
	structures := make(map[uint16][]packet.ValueDescriptor)
	for _, board := range boards {
		for _, packet := range board.Packets {
			id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
			if err != nil {
				return nil, err
			}
			structures[(uint16)(id)] = getDesciptor(packet.Values)
		}
	}
	return structures, nil
}

func getDesciptor(values []excel_models.Value) []packet.ValueDescriptor {
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
				enums[value.Name] = getEnumDescriptor(value.Type)
			}
		}
	}
	return enums
}

func getEnumDescriptor(literal string) packet.EnumDescriptor {
	return strings.Split(strings.ReplaceAll(strings.TrimSuffix(strings.TrimPrefix("enum(", literal), ")"), " ", ""), ",")
}

func getNames(global excel_models.GlobalInfo, boards map[string]excel_models.Board) (map[uint16][]string, error) {
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

// FIXME: remove hardcoded tags
func getMessageConfig(global excel_models.GlobalInfo, boards map[string]excel_models.Board) (message.Config, error) {
	fault, err := strconv.ParseUint(global.MessageToId["fault"], 10, 16)
	if err != nil {
		return message.Config{}, err
	}

	warning, err := strconv.ParseUint(global.MessageToId["warning"], 10, 16)
	if err != nil {
		return message.Config{}, err
	}

	blcuAck, err := strconv.ParseUint(global.MessageToId["blcu_ack"], 10, 16)
	if err != nil {
		return message.Config{}, err
	}

	return message.Config{
		FaultId:   (uint16)(fault),
		WarningId: (uint16)(warning),
		BlcuAckId: (uint16)(blcuAck),
	}, nil
}

func getIdKinds(global excel_models.GlobalInfo, boards map[string]excel_models.Board) (map[uint16]packet.Kind, error) {
	idToKind := make(map[uint16]packet.Kind)
	for _, board := range boards {
		for _, packet := range board.Packets {
			id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
			if err != nil {
				return nil, err
			}
			idToKind[(uint16)(id)] = getKind(packet.Description.Type)
		}
	}

	fault, err := strconv.ParseUint(global.MessageToId["fault"], 10, 16)
	if err != nil {
		return nil, err
	}

	warning, err := strconv.ParseUint(global.MessageToId["warning"], 10, 16)
	if err != nil {
		return nil, err
	}

	blcuAck, err := strconv.ParseUint(global.MessageToId["blcu_ack"], 10, 16)
	if err != nil {
		return nil, err
	}

	idToKind[(uint16)(fault)] = packet.Message
	idToKind[(uint16)(warning)] = packet.Message
	idToKind[(uint16)(blcuAck)] = packet.Message

	return idToKind, nil
}

func getKind(literal string) packet.Kind {
	switch literal {
	case "data":
		return packet.Data
	case "message":
		return packet.Message
	case "order":
		return packet.Order
	default:
		// TODO: handle error
		panic("unknown kind")
	}
}

func getFilter(addrs []string, protocolToPort map[string]string, tcpClientTag string, tcpServerTag string, udpTag string) string {
	// FIXME: IPIP filter
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
