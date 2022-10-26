package dto

import (
	podDataCreator "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/excelRetriever"
)

type BoardDTO struct {
	name         string
	descriptions map[Name]DescriptionDTO
	measurements map[Name]MeasurementDTO
	structures   map[Name]StructureDTO
}

func NewBoardDTO(sheet excelRetriever.Sheet) BoardDTO {
	return BoardDTO{
		name:         sheet.Name,
		descriptions: getDescriptionDTOs(sheet.Tables["Packet Description"]),
		measurements: getMeasurementDTOs(sheet.Tables["Value Description"]),
		structures:   getStructureDTOs(sheet.Tables["Packet Structure"]),
	}
}

func (b BoardDTO) GetPackets() map[Id]*podDataCreator.Packet {
	packetDTOs := b.GetPacketDTOs()
	packets := make(map[Id]*podDataCreator.Packet, len(packetDTOs))
	for _, packetDTO := range packetDTOs {
		packet := packetDTO.toPacket()
		packets[packet.Id] = &packet
	}

	return packets
}

func (b BoardDTO) GetPacketDTOs() []PacketDTO {
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

func getDescriptionDTOs(table excelRetriever.Table) map[Name]DescriptionDTO {
	descriptions := make(map[Name]DescriptionDTO, len(table.Rows))
	for _, row := range table.Rows {
		adapter := newDescriptionDTO(row)
		descriptions[adapter.Name] = adapter
	}

	return descriptions
}

func getMeasurementDTOs(table excelRetriever.Table) map[Name]MeasurementDTO {
	measurements := make(map[Name]MeasurementDTO, len(table.Rows))
	for _, row := range table.Rows {
		adapter := newMeasurementDTO(row)
		measurements[adapter.Name] = adapter
	}

	return measurements
}

func getStructureDTOs(table excelRetriever.Table) map[Name]StructureDTO {
	structures := make(map[Name]StructureDTO)
	for _, column := range getColumns(table) {
		structure := newStructureDTO(column)
		structures[structure.packetName] = structure
	}

	return structures
}

func getColumns(table excelRetriever.Table) [][]string {
	columns := make([][]string, len(table.Rows[0]))
	for i := 0; i < len(table.Rows[0]); i++ {
		columns[i] = getColumn(i, table)
	}

	return columns
}

func getColumn(i int, table excelRetriever.Table) []string {
	column := make([]string, len(table.Rows))
	for j := 0; j < len(table.Rows); j++ {
		column[j] = table.Rows[j][i]
	}

	return column
}
