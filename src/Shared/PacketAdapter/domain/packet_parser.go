package domain

import (
	"bytes"
	"io"
	"log"
	"strconv"

	order "github.com/HyperloopUPV-H8/Backend-H8/..."
	excelParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain/board"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain/serde"
)

type ID = uint16
type PacketMeasurements = map[string]MeasurementData
type Name = string

type PacketParser struct {
	packetTypes map[ID]PacketMeasurements
	enums       map[Name]Enum
}

func NewParser(packets []excelParser.Packet) PacketParser {
	return PacketParser{
		packetTypes: getPacketTypes(packets),
		enums:       getEnums(packets),
	}
}

func getEnums(packets []excelParser.Packet) map[Name]Enum {
	enums := make(map[Name]Enum, 0)
	for _, packet := range packets {
		for _, measurement := range packet.Measurements {
			if IsEnum(measurement.ValueType) {
				enums[measurement.Name] = NewEnum(measurement.ValueType)
			}
		}
	}

	return enums
}

func getPacketTypes(packets []excelParser.Packet) map[uint16]PacketMeasurements {
	packetMeasurements := make(map[ID]PacketMeasurements, len(packets))

	for _, packet := range packets {
		measurementDataArr := getMeasurementData(packet)
		id, err := strconv.ParseUint(packet.Description.ID, 10, 16)

		if err != nil {
			log.Fatal(err)
		}

		packetMeasurements[uint16(id)] = measurementDataArr
	}

	return packetMeasurements
}

func getMeasurementData(packet excelParser.Packet) map[string]MeasurementData {
	measurementDataArr := make(map[string]MeasurementData, len(packet.Measurements))
	for _, measurement := range packet.Measurements {
		valueType := measurement.ValueType
		if IsEnum(valueType) {
			valueType = "enum"
		}
		measurementData := NewMeasurement(measurement.Name, valueType)
		measurementDataArr[measurement.Name] = measurementData
	}
	return measurementDataArr
}

func (parser PacketParser) Decode(data []byte) PacketUpdate {
	dataReader := bytes.NewBuffer(data)
	id := serde.DecodeID(dataReader)
	values := parser.decodePacket(parser.packetTypes[id], dataReader)

	return NewPacketUpdate(id, values, dataReader.Bytes())
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

func (parser PacketParser) Encode(packet order.OrderWebAdapter) []byte {
	dataWriter := bytes.NewBuffer(make([]byte, 0))
	packetID, err := strconv.ParseUint(packet.ID, 10, 16)
	if err != nil {
		log.Fatalf("encode: %s\n", err)
	}
	serde.EncodeID(uint16(packetID), dataWriter)
	for name, value := range packet.Values {
		parser.encodeValue(uint16(packetID), name, value, dataWriter)
	}

	return dataWriter.Bytes()
}

func (parser PacketParser) encodeValue(id uint16, name string, value string, bytes io.Writer) {
	switch parser.packetTypes[id][name].valueType {
	case "enum":
		serde.EncodeEnum(parser.enums[name], value, bytes)
	case "bool":
		val, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatalf("encode value: %s\n", err)
		}
		serde.EncodeBool(val, bytes)
	default:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Fatalf("encode value: %s\n", err)
		}
		serde.EncodeNumber(parser.packetTypes[id][name].valueType, val, bytes)
	}
}
