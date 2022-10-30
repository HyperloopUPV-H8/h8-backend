package domain

import (
	"bytes"
	"io"
	"log"
	"strconv"

	parser "github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/application/interfaces"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain/interfaces"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/infra/serde"
)

type ID = uint16
type PacketMeasurements = []MeasurementData
type Name = string

type PacketParser struct {
	packetTypes map[ID]PacketMeasurements
	enums       map[Name]Enum
}

func NewParser(packets []parser.Packet) PacketParser {
	return PacketParser{
		packetTypes: getPacketTypes(packets),
		enums:       getEnums(packets),
	}
}

func getEnums(packets []parser.Packet) map[Name]Enum {
	enums := make(map[Name]Enum, 0)
	for _, packetDTO := range packets {
		for _, measurement := range packetDTO.Measurements() {
			if IsEnum(measurement.ValueType()) {
				enums[measurement.Name()] = NewEnum(measurement.ValueType())
			}
		}
	}

	return enums
}

func getPacketTypes(packets []parser.Packet) map[uint16]PacketMeasurements {
	packetMeasurements := make(map[ID]PacketMeasurements, len(packets))

	for _, packetDTO := range packets {
		measurementDataArr := getMeasurementData(packetDTO)
		id, err := strconv.ParseUint(packetDTO.Description().ID(), 10, 16)

		if err != nil {
			log.Fatal(err)
		}

		packetMeasurements[uint16(id)] = measurementDataArr
	}

	return packetMeasurements
}

func getMeasurementData(packet parser.Packet) []MeasurementData {
	measurementDataArr := make([]MeasurementData, len(packet.Measurements()))
	for index, measurementDTO := range packet.Measurements() {
		valueType := measurementDTO.ValueType()
		if IsEnum(valueType) {
			valueType = "enum"
		}
		measurementData := NewMeasurement(measurementDTO.Name(), valueType)
		measurementDataArr[index] = measurementData
	}
	return measurementDataArr
}

func (parser PacketParser) Decode(data []byte) interfaces.PacketUpdate {
	dataReader := bytes.NewBuffer(data)
	id := serde.DecodeID(dataReader)
	values := parser.decodePacket(parser.packetTypes[id], dataReader)

	return NewPacketUpdate(id, values)
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
		return serde.DecodeEnum(parser.enums[measurement.name], reader)
	case "bool":
		return serde.DecodeBool(reader)
	default:
		return serde.DecodeNumber(measurement.valueType, reader)
	}
}
