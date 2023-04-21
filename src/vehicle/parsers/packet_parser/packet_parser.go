package packet_parser

import (
	"bytes"
	"fmt"
	"io"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

type PacketParser struct {
	Ids          common.Set[uint16]
	structures   map[uint16][]packet.ValueDescriptor
	valueParsers map[string]parser
	config       Config
}

func NewPacketParser(ids common.Set[uint16], structures map[uint16][]packet.ValueDescriptor, enums map[string]packet.EnumDescriptor) PacketParser {
	return PacketParser{
		Ids:        ids,
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

func (parser *PacketParser) Decode(id uint16, raw []byte, metadata packet.Metadata) (models.PacketUpdate, error) {
	structure, ok := parser.structures[id]
	if !ok {
		return models.PacketUpdate{}, fmt.Errorf("structure for packet %d not found", id)
	}

	reader := bytes.NewReader(raw)

	values := make(map[string]packet.Value)
	for _, descriptor := range structure {
		value, err := parser.decodeValue(descriptor, reader)
		if err != nil {
			return models.PacketUpdate{}, err
		}

		values[descriptor.Name] = value
	}

	return models.PacketUpdate{
		Metadata: metadata,
		HexValue: raw,
		Values:   values,
	}, nil
}

func (parser *PacketParser) decodeValue(descriptor packet.ValueDescriptor, data io.Reader) (packet.Value, error) {
	decoder, ok := parser.valueParsers[descriptor.Type]
	if !ok {
		return nil, fmt.Errorf("decoder for type %s not found", descriptor.Type)
	}

	return decoder.decode(descriptor, parser.config.GetByteOrder(), data)
}

func (parser *PacketParser) Encode(id uint16, values map[string]packet.Value, data io.Writer) error {
	structure, ok := parser.structures[id]
	if !ok {
		return fmt.Errorf("structure for packet %d not found", id)
	}

	for _, descriptor := range structure {
		value, ok := values[descriptor.Name]
		if !ok {
			return fmt.Errorf("value for %s not found", descriptor.Name)
		}

		err := parser.encodeValue(descriptor, value, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (parser *PacketParser) encodeValue(descriptor packet.ValueDescriptor, value packet.Value, data io.Writer) error {
	encoder, ok := parser.valueParsers[descriptor.Type]
	if !ok {
		return fmt.Errorf("encoder for type %s not found", descriptor.Type)
	}

	return encoder.encode(descriptor, parser.config.GetByteOrder(), value, data)
}
