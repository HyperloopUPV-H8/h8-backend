package infra

import (
	"bytes"
	"log"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra/dto"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra/internals/decoder"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra/internals/encoder"
)

func Decode(data []byte, aggregate PacketAggregate) dto.PacketUpdate {
	reader := bytes.NewBuffer(data)
	id := decoder.ID(reader)

	packet := aggregate.get(id)
	if packet == nil {
		log.Fatalf("packet parser: decode: invalid ID %d\n", id)
	}

	values := make(map[string]any, len(packet))
	for _, value := range packet {
		switch value.Kind {
		case "enum":
			values[value.Name] = decoder.Enum(reader, aggregate.getEnum(value.Name))
		case "bool":
			values[value.Name] = decoder.Bool(reader)
		case "string":
			values[value.Name] = decoder.String(reader)
		default:
			values[value.Name] = decoder.Number(reader, value.Kind)
		}
	}

	return dto.NewPacketUpdate(id, values, data)
}

func Encode(data dto.PacketUpdate, aggregate PacketAggregate) []byte {
	writer := bytes.NewBuffer([]byte{})
	id := data.ID()
	encoder.ID(writer, id)

	packetValues := data.Values()
	for _, value := range aggregate.get(id) {
		switch value.Kind {
		case "enum":
			encoder.Enum(writer, aggregate.getEnum(value.Name), packetValues[value.Name].(string))
		case "bool":
			encoder.Bool(writer, packetValues[value.Name].(bool))
		case "string":
			encoder.String(writer, packetValues[value.Name].(string))
		default:
			encoder.Number(writer, packetValues[value.Name].(float64), value.Kind)
		}
	}

	return writer.Bytes()
}
