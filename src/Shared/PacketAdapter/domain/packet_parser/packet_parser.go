package packetparser

import (
	"bytes"
	"io"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain/packet_parser/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain/packet_parser/infra"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/excelAdapter/dto"
)

type PacketParser struct {
	packetTypes map[domain.ID]domain.PacketMeasurements
	enums       map[domain.Name]domain.Enum
}

func New(packetDTOs []dto.PacketDTO) PacketParser {
	return PacketParser{} //FIXME:
}

func (parser PacketParser) Decode(data []byte) domain.PacketUpdate {
	dataReader := bytes.NewBuffer(data)
	id := infra.DecodeID(dataReader)

	values := parser.decodePacket(parser.packetTypes[id], dataReader)

	return domain.NewUpdatedValues(id, values)
}

func (parser PacketParser) decodePacket(measurements domain.PacketMeasurements, bytes io.Reader) map[domain.Name]any {
	values := make(map[domain.Name]any, len(measurements))
	for _, measure := range measurements {
		values[measure.Name] = parser.decodeMeasurement(measure, bytes)
	}
	return values
}

func (parser PacketParser) decodeMeasurement(measurement domain.MeasurementData, reader io.Reader) any {
	switch measurement.ValueType {
	case "enum":
		return infra.DecodeEnum(parser.enums[measurement.Name], reader)
	case "bool":
		return infra.DecodeBool(reader)
	default:
		return infra.DecodeNumber(measurement.ValueType, reader)
	}
}
