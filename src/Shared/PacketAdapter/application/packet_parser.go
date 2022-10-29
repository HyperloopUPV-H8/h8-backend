package application

import (
	"bytes"
	"io"
	"log"
	"strconv"

	parser "github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/application/interfaces"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/application/interfaces"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/infra/serde"
)

type ID = uint16
type PacketMeasurements = []interfaces.Measurement
type Name = string
type EnumVariant = string
type Enum = map[uint8]EnumVariant

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
			if domain.IsEnum(measurement.ValueType()) {
				enums[measurement.Name()] = domain.NewEnum(measurement.ValueType())
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

func getMeasurementData(packet parser.Packet) []interfaces.Measurement {
	measurementDataArr := make([]interfaces.Measurement, len(packet.Measurements()))
	for index, measurementDTO := range packet.Measurements() {
		valueType := measurementDTO.ValueType()
		if domain.IsEnum(valueType) {
			valueType = "enum"
		}
		measurementData := domain.NewMeasurement(measurementDTO.Name(), valueType)
		measurementDataArr[index] = measurementData
	}
	return measurementDataArr
}

func (parser PacketParser) Decode(data []byte) interfaces.PacketUpdate {
	dataReader := bytes.NewBuffer(data)
	id := serde.DecodeID(dataReader)
	values := parser.decodePacket(parser.packetTypes[id], dataReader)

	return domain.NewPacketUpdate(id, values)
}

func (parser PacketParser) decodePacket(measurements PacketMeasurements, bytes io.Reader) map[Name]any {
	values := make(map[Name]any, len(measurements))
	for _, measure := range measurements {
		values[measure.Name()] = parser.decodeMeasurement(measure, bytes)
	}
	return values
}

func (parser PacketParser) decodeMeasurement(measurement interfaces.Measurement, reader io.Reader) any {
	switch measurement.ValueType() {
	case "enum":
		return serde.DecodeEnum(parser.enums[measurement.Name()], reader)
	case "bool":
		return serde.DecodeBool(reader)
	default:
		return serde.DecodeNumber(measurement.ValueType(), reader)
	}
}
