package domain

import (
	"bytes"
	"io"

	exceladapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excelAdapter/dto"
)

type PacketParser struct {
	packetTypes map[ID]PacketMeasurements
	enums       map[Name]Enum
}

func New(packetDTOs []exceladapter.PacketDTO) PacketParser {

}

func (parser PacketParser) Decode(data []byte) PacketUpdate {
	dataReader := bytes.NewBuffer(data)
	id := DecodeID(dataReader)

	values := parser.decodePacket(parser.packetTypes[id], dataReader)

	return NewUpdatedValues(id, values)
}

func (parser PacketParser) decodePacket(measurements PacketMeasurements, bytes io.Reader) map[Name]any {
	values := make(map[Name]any, len(measurements))
	for _, measure := range measurements {
		values[measure.name] = parser.decodeMeasurement(measure, bytes)
	}
	return values
}

func (parser PacketParser) decodeMeasurement(measurement MeasurementData, reader io.Reader) any {
	switch measurement.valueType {
	case "enum":
		return DecodeEnum(parser.enums[measurement.name], reader)
	case "bool":
		return DecodeBool(reader)
	default:
		return DecodeNumber(measurement.valueType, reader)
	}
}
