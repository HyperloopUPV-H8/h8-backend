package packet_parser

import (
	"bytes"
	"io"
	"strconv"
	"strings"

	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/packet_parser/internals"
	"github.com/HyperloopUPV-H8/Backend-H8/packet_parser/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

type PacketParser struct {
	descriptors map[uint16]models.PacketDescriptor
	enums       map[string]models.Enum
	trace       zerolog.Logger
}

func New(boards map[string]excelAdapterModels.Board) PacketParser {
	trace.Info().Msg("new packet parser")

	return PacketParser{
		descriptors: getDescriptors(boards),
		enums:       getEnums(boards),
		trace:       trace.With().Str("component", "packetParser").Logger(),
	}
}

func getDescriptors(boards map[string]excelAdapterModels.Board) map[uint16]models.PacketDescriptor {
	descriptors := make(map[uint16]models.PacketDescriptor)

	for _, board := range boards {
		for _, packet := range board.Packets {
			id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
			if err != nil {
				continue
			}
			descriptors[uint16(id)] = getValueDescriptors(packet.Values)
		}
	}

	return descriptors
}

func getValueDescriptors(values []excelAdapterModels.Value) []models.ValueDescriptor {
	valueDescriptors := make([]models.ValueDescriptor, 0, len(values))
	enums := make(map[string]models.Enum)
	for _, value := range values {
		if value.ID == "" {
			continue
		}

		kind := value.Type
		if strings.HasPrefix(strings.ToUpper(kind), "ENUM") {
			kind = "enum"
			enums[value.ID] = models.GetEnum(strings.ToUpper(value.Type))
		}

		valueDescriptors = append(valueDescriptors, models.ValueDescriptor{
			ID:   value.ID,
			Type: kind,
		})
	}

	return valueDescriptors
}

func getEnums(boards map[string]excelAdapterModels.Board) map[string]models.Enum {
	enums := make(map[string]models.Enum)

	for _, board := range boards {
		for _, packet := range board.Packets {
			for _, value := range packet.Values {
				if value.ID == "" {
					continue
				}

				kind := value.Type
				if strings.HasPrefix(strings.ToUpper(kind), "ENUM") {
					kind = "enum"
					enums[value.ID] = models.GetEnum(strings.ToUpper(value.Type))
				}
			}
		}
	}

	return enums
}

func (parser PacketParser) Decode(raw []byte) (id uint16, values models.PacketValues) {
	reader := bytes.NewReader(raw)
	id = internals.DecodeID(reader)

	parser.trace.Trace().Uint16("id", id).Msg("decode")
	values = make(models.PacketValues, len(parser.descriptors[id]))
	for _, value := range parser.descriptors[id] {
		values[value.ID] = parser.decodeValue(value, reader)
	}

	return id, values
}

func (parser PacketParser) decodeValue(value models.ValueDescriptor, reader io.Reader) any {
	parser.trace.Trace().Str("type", value.Type).Msg("decode value")

	switch value.Type {
	case "enum":
		value := internals.DecodeEnum(reader, parser.enums[value.ID])
		return value
	case "bool":
		return internals.DecodeBool(reader)
	case "string":
		return internals.DecodeString(reader)
	default:
		return internals.DecodeNumber(reader, value.Type)
	}
}

func (parser PacketParser) Encode(id uint16, values models.PacketValues) []byte {
	parser.trace.Trace().Uint16("id", id).Msg("encode")

	writer := bytes.NewBuffer([]byte{})
	internals.EncodeID(writer, id)

	for _, valueDescriptor := range parser.descriptors[id] {
		parser.encodeValue(valueDescriptor, values[valueDescriptor.ID], writer)
	}

	return writer.Bytes()
}

func (parser PacketParser) encodeValue(desc models.ValueDescriptor, value any, writer io.Writer) {
	parser.trace.Trace().Str("type", desc.Type).Msg("encode value")

	switch desc.Type {
	case "enum":
		internals.EncodeEnum(writer, parser.enums[desc.ID], value.(string))
	case "bool":
		internals.EncodeBool(writer, value.(bool))
	case "string":
		internals.EncodeString(writer, value.(string))
	default:
		internals.EncodeNumber(writer, desc.Type, value.(float64))
	}
}
