package domain

import (
	"bytes"
	"io"

	value "github.com/HyperloopUPV-H8/Backend-H8/..."
)

type PacketParser struct {
	packetTypes map[ID]PacketMeasurements
	enums       map[Name]Enum
}

func (parser PacketParser) Decode(data []byte) UpdatedValues {
	dataReader := bytes.NewBuffer(data)
	id := DecodeID(dataReader)

	values := parser.decodePacket(parser.packetTypes[id], dataReader)

	return NewUpdatedValues(id, values)
}

func (parser PacketParser) decodePacket(measurements PacketMeasurements, bytes io.Reader) map[Name]value.Value {
	values := make(map[Name]value.Value, len(measurements))
	for _, measure := range measurements {
		values[measure.name] = parser.decodeMeasurement(measure, bytes)
	}
	return values
}

func (parser PacketParser) decodeMeasurement(measurement MeasurementData, reader io.Reader) value.Value {
	switch measurement.valueType {
	case "enum":
		return DecodeEnum(parser.enums[measurement.name], reader)
	case "bool":
		return DecodeBool(reader)
	default:
		return DecodeNumber(measurement.valueType, reader)
	}
}
