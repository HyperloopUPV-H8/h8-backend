package vehicle

import (
	"strconv"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	excel_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/HyperloopUPV-H8/Backend-H8/pipe"
	"github.com/HyperloopUPV-H8/Backend-H8/sniffer"
	"github.com/HyperloopUPV-H8/Backend-H8/unit_converter"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/packet_parser"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/protection_parser"
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

	packetParser, err := packet_parser.CreatePacketParser(args.GlobalInfo, args.Boards, vehicleTrace)

	if err != nil {
		vehicleTrace.Fatal().Err(err).Msg("error creating packetParser")
	}

	names, err := getPacketToValuesNames(args.GlobalInfo, args.Boards)

	if err != nil {
		vehicleTrace.Error().Err(err).Msg("error getting packet to values names")
	}

	snifferConfig := getSnifferConfig(args.Config)
	pipesConfig := getPipesConfig(args.Config)

	protectionIds := common.NewSet[uint16]()

	faultId := mustGetId(args.GlobalInfo.MessageToId, args.Config.Protections.FaultIdKey, vehicleTrace)
	protectionIds.Add(faultId)
	warningId := mustGetId(args.GlobalInfo.MessageToId, args.Config.Protections.WarningIdKey, vehicleTrace)
	protectionIds.Add(warningId)
	errorId := mustGetId(args.GlobalInfo.MessageToId, args.Config.Protections.ErrorIdKey, vehicleTrace)
	protectionIds.Add(errorId)

	vehicle := Vehicle{
		podConverter:     unit_converter.NewUnitConverter("pod", args.Boards, args.GlobalInfo.UnitToOperations),
		displayConverter: unit_converter.NewUnitConverter("display", args.Boards, args.GlobalInfo.UnitToOperations),

		sniffer: sniffer.CreateSniffer(args.GlobalInfo, snifferConfig, vehicleTrace),
		pipes:   pipe.CreatePipes(args.GlobalInfo, dataChan, args.OnConnectionChange, pipesConfig, vehicleTrace),

		dataIds:       getBoardIdsFromType(args.Boards, "data"),
		orderIds:      getBoardIdsFromType(args.Boards, "order"),
		protectionIds: protectionIds,

		packetParser:     packetParser,
		protectionParser: protection_parser.NewProtectionParser(args.GlobalInfo, faultId, warningId, errorId),
		bitarrayParser:   NewBitarrayParser(names),

		dataChan: dataChan,

		idToBoard:          getIdToBoard(args.Boards, vehicleTrace),
		onConnectionChange: args.OnConnectionChange,
		trace:              vehicleTrace,
	}

	vehicle.sniffer.Listen(dataChan)

	return vehicle
}

func getSnifferConfig(config Config) sniffer.Config {
	return sniffer.Config{
		TcpClientTag: config.Network.TcpClientTag,
		TcpServerTag: config.Network.TcpServerTag,
		UdpTag:       config.Network.UdpTag,
		Mtu:          config.Network.Mtu,
		Interface:    config.Network.Interface,
	}
}

func getPipesConfig(config Config) pipe.Config {
	return pipe.Config{
		TcpClientTag: config.Network.TcpClientTag,
		TcpServerTag: config.Network.TcpServerTag,
		Mtu:          config.Network.Mtu,
	}
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

func getBoardIdsFromType(boards map[string]excel_models.Board, kind string) common.Set[uint16] {
	ids := common.NewSet[uint16]()

	for _, board := range boards {
		for _, packet := range board.Packets {
			if packet.Description.Type == kind {
				id, err := strconv.ParseInt(packet.Description.ID, 10, 16)

				if err != nil {
					//TODO:
				}

				ids.Add(uint16(id))
			}
		}
	}

	return ids
}

func mustGetId(kindToId map[string]string, key string, trace zerolog.Logger) uint16 {
	idStr, ok := kindToId[key]

	if !ok {
		trace.Fatal().Str("key", key).Msg("key not found")
	}

	id, err := strconv.ParseUint(idStr, 10, 16)

	if err != nil {
		trace.Fatal().Str("id", idStr).Msg("error parsing id")
	}

	return uint16(id)
}
