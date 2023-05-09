package packet_parser

import (
	"strconv"
	"strings"

	excel_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/rs/zerolog"
)

func CreatePacketParser(global excel_models.GlobalInfo, boards map[string]excel_models.Board, trace zerolog.Logger) (PacketParser, error) {
	structures, err := getStructures(global, boards)
	if err != nil {
		return PacketParser{}, err
	}

	return newPacketParser(structures, getEnumDescriptors(global, boards)), nil
}

func newPacketParser(structures map[uint16][]packet.ValueDescriptor, enums map[string]packet.EnumDescriptor) PacketParser {
	return PacketParser{
		structures: structures,
		valueParsers: map[string]parser{
			"uint8":   numericParser[uint8]{},
			"uint16":  numericParser[uint16]{},
			"uint32":  numericParser[uint32]{},
			"uint64":  numericParser[uint64]{},
			"int8":    numericParser[int8]{},
			"int16":   numericParser[int16]{},
			"int32":   numericParser[int32]{},
			"int64":   numericParser[int64]{},
			"float32": numericParser[float32]{},
			"float64": numericParser[float64]{},
			"bool":    booleanParser{},
			"enum":    enumParser{descriptors: enums},
		},
	}
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

func getEnumDescriptors(global excel_models.GlobalInfo, boards map[string]excel_models.Board) map[string]packet.EnumDescriptor {
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
