package vehicle

import (
	"strconv"

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

	vehicle := Vehicle{
		podConverter:     unit_converter.NewUnitConverter("pod", args.Boards, args.GlobalInfo.UnitToOperations),
		displayConverter: unit_converter.NewUnitConverter("display", args.Boards, args.GlobalInfo.UnitToOperations),

		sniffer: sniffer.CreateSniffer(args.GlobalInfo, snifferConfig, vehicleTrace),
		pipes:   pipe.CreatePipes(args.GlobalInfo, dataChan, args.OnConnectionChange, pipesConfig, vehicleTrace),

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
