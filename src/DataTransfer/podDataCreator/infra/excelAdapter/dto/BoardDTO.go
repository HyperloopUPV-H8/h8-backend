package dto

import (
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/excelRetreiver"
	podDataCreator "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain"
)

type BoardDTO struct {
	descriptions map[string]DescriptionDTO
	measurements map[string]MeasurementDTO
	structures   map[string]StructureDTO
}

func NewBoardDTO(tables map[string]excelRetreiver.Table) BoardDTO {
	measurements := getMeasurementDTOs(tables["ValueDescription"])
	descriptions := getDescriptionDTOs(tables["PacketDescription"])
	structures := getStructureDTOs(tables["PacketStructure"])

	return BoardDTO{
		descriptions: descriptions,
		measurements: measurements,
		structures:   structures,
	}
}

func (b *BoardDTO) GetPackets() map[uint16]*podDataCreator.Packet {
	packetDTOs := b.getPacketDTOs()
	packets := make(map[uint16]*podDataCreator.Packet, len(packetDTOs))
	for _, packetDTO := range packetDTOs {
		packet := packetDTO.toPacket()
		packets[packet.Id] = &packet
	}

	return packets
}

func (b *BoardDTO) getPacketDTOs() []PacketDTO {
	expandedPacketDTOs := make([]PacketDTO, 0)
	for _, description := range b.descriptions {
		measurementDTOs := b.getPacketMeasurements(description)
		packetDTOs := expandPacketDTO(description, measurementDTOs)
		expandedPacketDTOs = append(expandedPacketDTOs, packetDTOs...)
	}
	return expandedPacketDTOs
}

func (b *BoardDTO) getPacketMeasurements(description DescriptionDTO) []MeasurementDTO {
	measurements := make([]MeasurementDTO, 0)

	for _, name := range b.structures[description.Name].measurements {
		measurements = append(measurements, b.measurements[name])
	}

	return measurements
}

func getDescriptionDTOs(table excelRetreiver.Table) map[string]DescriptionDTO {
	descriptions := make(map[string]DescriptionDTO)
	for _, row := range table.Rows {
		adapter := newDescriptionDTO(row)
		descriptions[adapter.Name] = adapter
	}

	return descriptions
}

func getMeasurementDTOs(table excelRetreiver.Table) map[string]MeasurementDTO {
	measurements := make(map[string]MeasurementDTO)
	for _, row := range table.Rows {
		adapter := newMeasurementDTO(row)
		measurements[adapter.Name] = adapter
	}

	return measurements
}

func getStructureDTOs(table excelRetreiver.Table) map[string]StructureDTO {
	columns := make([][]string, 0)
	for i := 0; i < len(table.Rows[0]); i++ {
		column := make([]string, 0)
		for j := 0; j < len(table.Rows); j++ {
			column = append(column, table.Rows[j][i])
		}
		columns = append(columns, column)
	}

	structures := make(map[string]StructureDTO)

	for _, column := range columns {
		structure := newStructure(column)
		structures[structure.packetName] = structure
	}

	return structures
}
