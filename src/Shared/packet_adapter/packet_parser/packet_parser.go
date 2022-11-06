package domain

import (
	"bytes"
	"io"
	"log"
	"strconv"

	ordertransfer "github.com/HyperloopUPV-H8/Backend-H8/OrderTransfer"
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra/serde"
)

type ID = uint16
type PacketMeasurements = []domain.MeasurementData
type Name = string

type PacketParser struct {
	packetTypes map[ID]PacketMeasurements
	enums       map[Name]domain.Enum
}

func NewParser(packets []excelAdapter.PacketDTO) PacketParser {
	return PacketParser{
		packetTypes: getPacketTypes(packets),
		enums:       getEnums(packets),
	}
}

func getEnums(packets []excelAdapter.PacketDTO) map[Name]domain.Enum {
	enums := make(map[Name]domain.Enum, 0)
	for _, packet := range packets {
		for _, measurement := range packet.Measurements {
			if domain.IsEnum(measurement.ValueType) {
				enums[measurement.Name] = domain.NewEnum(measurement.ValueType)
			}
		}
	}

	return enums
}

func getPacketTypes(packets []excelAdapter.PacketDTO) map[uint16]PacketMeasurements {
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

func getMeasurementData(packet excelAdapter.PacketDTO) []domain.MeasurementData {
	measurementDataArr := make([]domain.MeasurementData, len(packet.Measurements))
	for index, measurement := range packet.Measurements {
		valueType := measurement.ValueType
		if domain.IsEnum(valueType) {
			valueType = "enum"
		}
		measurementData := domain.NewMeasurement(measurement.Name, valueType)
		measurementDataArr[index] = measurementData
	}
	return measurementDataArr
}

func (parser PacketParser) Decode(data []byte) domain.PacketUpdate {
	dataReader := bytes.NewBuffer(data)
	id := serde.DecodeID(dataReader)
	values := parser.decodePacket(parser.packetTypes[id], dataReader)

	return domain.NewPacketUpdate(id, values, dataReader.Bytes())
}

func (parser PacketParser) decodePacket(measurements PacketMeasurements, bytes io.Reader) map[Name]any {
	values := make(map[Name]any, len(measurements))
	for _, measurementData := range measurements {
		values[measurementData.Name] = parser.decodeMeasurement(measurementData, bytes)
	}
	return values
}

func (parser PacketParser) decodeMeasurement(measurement domain.MeasurementData, reader io.Reader) any {
	switch measurement.ValueType {
	case "enum":
		return serde.DecodeEnum(parser.enums[measurement.Name], reader)
	case "bool":
		return serde.DecodeBool(reader)
	case "string":
		return serde.DecodeString(reader)
	default:
		return serde.DecodeNumber(measurement.ValueType, reader)
	}
}

func (parser PacketParser) Encode(packet ordertransfer.OrderWebAdapter) []byte {
	dataWriter := bytes.NewBuffer(make([]byte, 0))
	serde.EncodeID(packet.Id, dataWriter)
	for name, value := range packet.Fields {
		parser.encodeValue(packet.Id, name, value, dataWriter)
	}

	return dataWriter.Bytes()
}

func (parser PacketParser) encodeValue(id uint16, name string, value string, bytes io.Writer) {
	switch parser.findMeasurement(id, name).ValueType {
	case "enum":
		serde.EncodeEnum(parser.enums[name], value, bytes)
	case "bool":
		val, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatalf("encode value: %s\n", err)
		}
		serde.EncodeBool(val, bytes)
	case "string":
		serde.EncodeString(value, bytes)
	default:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Fatalf("encode value: %s\n", err)
		}
		serde.EncodeNumber(parser.findMeasurement(id, name).ValueType, val, bytes)
	}
}

func (parser PacketParser) findMeasurement(id uint16, name string) domain.MeasurementData {
	for _, measurement := range parser.packetTypes[id] {
		if measurement.Name == name {
			return measurement
		}
	}
	panic("measurement not found")
}
