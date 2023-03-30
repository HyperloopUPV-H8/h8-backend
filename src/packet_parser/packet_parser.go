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

func (parser *PacketParser) AddGlobal(excelAdapterModels.GlobalInfo) {
	parser.trace.Debug().Msg("add global")
}

func (parser *PacketParser) AddPacket(boardName string, packet excelAdapterModels.Packet) {
	parser.trace.Debug().Str("id", packet.Description.ID).Str("name", packet.Description.Name).Str("board", boardName).Msg("add packet")

	id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
	if err != nil {
		parser.trace.Error().Err(err).Str("id", packet.Description.ID).Msg("")
		return
	}

	valueDescriptors := make([]models.ValueDescriptor, 0, len(packet.Values))
	for _, value := range packet.Values {
		if value.ID == "" {
			continue
		}

		parser.trace.Trace().Str("id", value.ID).Str("type", value.Type).Msg("add value")

		kind := value.Type
		if strings.HasPrefix(strings.ToUpper(kind), "ENUM") {
			kind = "enum"
			parser.enums[value.ID] = models.GetEnum(strings.ToUpper(value.Type))
		}

		valueDescriptors = append(valueDescriptors, models.ValueDescriptor{
			ID:   value.ID,
			Type: kind,
		})
	}

	parser.descriptors[uint16(id)] = valueDescriptors
}

func NewPacketParser() *PacketParser {
	trace.Info().Msg("new packet parser")
	return &PacketParser{
		descriptors: make(map[uint16]models.PacketDescriptor),
		enums:       make(map[string]models.Enum),
		trace:       trace.With().Str("component", "packetParser").Logger(),
	}
}

func (parser PacketParser) Decode(raw []byte) (id uint16, values models.PacketValues) {
	reader := bytes.NewReader(raw)
	id = internals.DecodeID(reader)

	values = make(models.PacketValues, len(parser.descriptors[id]))
	for _, value := range parser.descriptors[id] {
		values[value.ID] = parser.decodeValue(value, reader)
	}

	parser.trace.Trace().Uint16("id", id).Msg("decode")
	return id, values
}

func (parser PacketParser) decodeValue(value models.ValueDescriptor, reader io.Reader) any {
	parser.trace.Trace().Str("type", value.Type).Msg("decode value")

	switch value.Type {
	case "enum":
		return internals.DecodeEnum(reader, parser.enums[value.ID])
	case "bool":
		return internals.DecodeBool(reader)
	case "string":
		return internals.DecodeString(reader)
	default:
		return internals.DecodeNumber(reader, value.Type)
	}
}

func (parser PacketParser) Encode(id uint16, values models.PacketValues) []byte {
	writer := bytes.NewBuffer([]byte{})
	internals.EncodeID(writer, id)

	for _, valueDescriptor := range parser.descriptors[id] {
		parser.encodeValue(valueDescriptor, values[valueDescriptor.ID], writer)
	}

	parser.trace.Trace().Uint16("id", id).Msg("encode")
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
