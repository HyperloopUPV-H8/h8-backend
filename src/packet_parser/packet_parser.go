package packet_parser

import (
	"bytes"
	"io"
	"log"
	"strconv"
	"strings"

	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/packet_parser/internals"
	"github.com/HyperloopUPV-H8/Backend-H8/packet_parser/models"
)

type PacketParser struct {
	descriptors map[uint16]models.PacketDescriptor
	enums       map[string]models.Enum
}

func (parser *PacketParser) AddPacket(board string, ip string, desc excelAdapterModels.Description, values []excelAdapterModels.Value) {
	id, err := strconv.ParseUint(desc.ID, 10, 16)
	if err != nil {
		log.Fatalf("packet parser: AddPacket: %s\n", err)
	}

	valueDescriptors := make([]models.ValueDescriptor, len(values))
	for i, value := range values {
		kind := value.Type
		if strings.HasPrefix(kind, "ENUM") {
			kind = "enum"
			parser.enums[value.Name] = models.GetEnum(value.Type)
		}

		valueDescriptors[i] = models.ValueDescriptor{
			Name: value.Name,
			Type: kind,
		}
	}

	parser.descriptors[uint16(id)] = valueDescriptors
}

func NewPacketParser() PacketParser {
	return PacketParser{
		descriptors: make(map[uint16]models.PacketDescriptor),
		enums:       make(map[string]models.Enum),
	}
}

func (parser PacketParser) Decode(raw []byte) (id uint16, values models.PacketValues) {
	reader := bytes.NewReader(raw)
	id = internals.DecodeID(reader)

	values = make(models.PacketValues, len(parser.descriptors[id]))
	for _, value := range parser.descriptors[id] {
		values[value.Name] = parser.decodeValue(value, reader)
	}

	return id, values
}

func (parser PacketParser) decodeValue(value models.ValueDescriptor, reader io.Reader) any {
	switch value.Type {
	case "enum":
		return internals.DecodeEnum(reader, parser.enums[value.Name])
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
		parser.encodeValue(valueDescriptor, values[valueDescriptor.Name], writer)
	}
	return writer.Bytes()
}

func (parser PacketParser) encodeValue(desc models.ValueDescriptor, value any, writer io.Writer) {
	switch desc.Type {
	case "enum":
		internals.EncodeEnum(writer, parser.enums[desc.Name], value.(string))
	case "bool":
		internals.EncodeBool(writer, value.(bool))
	case "string":
		internals.EncodeString(writer, value.(string))
	default:
		internals.EncodeNumber(writer, desc.Type, value.(float64))
	}
}
