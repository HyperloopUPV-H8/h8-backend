package packetparser

import (
	"bytes"
	"io"
	"log"
	"strconv"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain/packet_parser/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain/packet_parser/infra"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/excelAdapter/dto"
)

type PacketParser struct {
	packetTypes map[domain.ID]domain.PacketMeasurements
	enums       map[domain.Name]domain.Enum
}

func New(packetDTOs []dto.PacketDTO) PacketParser {

	return PacketParser{
		packetTypes: getPacketTypes(packetDTOs),
		enums:       getEnums(packetDTOs),
	}
}

func getEnums(packetDTOs []dto.PacketDTO) map[domain.Name]domain.Enum {
	enums := make(map[domain.Name]domain.Enum, 0)
	for _, packetDTO := range packetDTOs {
		for _, measurement := range packetDTO.Measurements {
			if domain.IsEnum(measurement.ValueType) {
				enums[measurement.Name] = domain.NewEnum(measurement.ValueType)
			}
		}
	}

	return enums
}

func getPacketTypes(packetDTOs []dto.PacketDTO) map[uint16]domain.PacketMeasurements {
	packetMeasurements := make(map[domain.ID]domain.PacketMeasurements, len(packetDTOs))

	for _, packetDTO := range packetDTOs {
		measurementDataArr := getMeasurementData(packetDTO)
		id, err := strconv.Atoi(packetDTO.Description.Id)

		if err != nil {
			log.Fatal(err)
		}

		packetMeasurements[uint16(id)] = measurementDataArr
	}

	return packetMeasurements
}

func getMeasurementData(packetDTO dto.PacketDTO) []domain.MeasurementData {
	measurementDataArr := make([]domain.MeasurementData, len(packetDTO.Measurements))
	for index, measurementDTO := range packetDTO.Measurements {
		valueType := measurementDTO.ValueType
		if domain.IsEnum(valueType) {
			valueType = "ENUM"
		}
		measurementData := domain.MeasurementData{
			Name:      measurementDTO.Name,
			ValueType: valueType,
		}
		measurementDataArr[index] = measurementData
	}
	return measurementDataArr
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
