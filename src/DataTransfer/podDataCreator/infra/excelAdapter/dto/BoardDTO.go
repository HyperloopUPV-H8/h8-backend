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
	return BoardDTO{
		descriptions: getDescriptionDTOs(tables["PacketDescription"]),
		measurements: getMeasurementDTOs(tables["ValueDescription"]),
		structures:   getStructureDTOs(tables["PacketStructure"]),
	}
}

func (b BoardDTO) GetPackets() map[uint16]*podDataCreator.Packet {
	packetDTOs := b.getPacketDTOs()
	packets := make(map[uint16]*podDataCreator.Packet, len(packetDTOs))
	for _, packetDTO := range packetDTOs {
		packet := packetDTO.toPacket()
		packets[packet.Id] = &packet
	}

	return packets
}

func (b BoardDTO) getPacketDTOs() []PacketDTO {
	expandedPacketDTOs := make([]PacketDTO, 0)
	for _, description := range b.descriptions {
		measurementDTOs := b.getPacketMeasurements(description)
		packetDTOs := expandPacketDTO(description, measurementDTOs)
		expandedPacketDTOs = append(expandedPacketDTOs, packetDTOs...)
	}
	return expandedPacketDTOs
}

func (b BoardDTO) getPacketMeasurements(description DescriptionDTO) []MeasurementDTO {
	wantedMeasurements := b.structures[description.Name].measurements
	measurements := make([]MeasurementDTO, len(wantedMeasurements))
	for index, name := range wantedMeasurements {
		measurements[index] = b.measurements[name]
	}

	return measurements
}

func getDescriptionDTOs(table excelRetreiver.Table) map[string]DescriptionDTO {
	descriptions := make(map[string]DescriptionDTO, len(table.Rows))
	for _, row := range table.Rows {
		adapter := newDescriptionDTO(row)
		descriptions[adapter.Name] = adapter
	}

	return descriptions
}

func getMeasurementDTOs(table excelRetreiver.Table) map[string]MeasurementDTO {
	measurements := make(map[string]MeasurementDTO, len(table.Rows))
	for _, row := range table.Rows {
		adapter := newMeasurementDTO(row)
		measurements[adapter.Name] = adapter
	}

	return measurements
}

func getStructureDTOs(table excelRetreiver.Table) map[string]StructureDTO {
	structures := make(map[string]StructureDTO)
	for _, column := range getColumns(table) {
		structure := newStructureDTO(column)
		structures[structure.packetName] = structure
	}

	return structures
}

func getColumns(table excelRetreiver.Table) [][]string {
	columns := make([][]string, len(table.Rows[0]))
	for i := 0; i < len(table.Rows[0]); i++ {
		column := make([]string, len(table.Rows))
		for j := 0; j < len(table.Rows); j++ {
			column[j] = table.Rows[j][i]
		}
		columns[i] = column
	}

	return columns
}