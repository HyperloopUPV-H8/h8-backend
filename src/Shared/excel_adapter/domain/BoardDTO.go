package domain

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_retriever/domain"
)

type BoardDTO struct {
	Name         string
	Descriptions map[Name]DescriptionDTO
	Measurements map[Name]MeasurementDTO
	Structures   map[Name]Structure
}

func NewBoard(sheet domain.Sheet) BoardDTO {
	return BoardDTO{
		Name:         sheet.Name,
		Descriptions: getDescriptions(sheet.Tables["PacketDescription"]),
		Measurements: getMeasurements(sheet.Tables["ValueDescription"]),
		Structures:   getStructures(sheet.Tables["PacketStructure"]),
	}
}

func (board BoardDTO) GetPackets() []PacketDTO {
	expandedPackets := make([]PacketDTO, 0)
	for _, description := range board.Descriptions {
		measurements := board.getPacketMeasurements(description)
		packetDTOs := expandPacket(description, measurements)
		expandedPackets = append(expandedPackets, packetDTOs...)
	}
	return expandedPackets
}

func (board BoardDTO) getPacketMeasurements(description DescriptionDTO) []MeasurementDTO {
	wantedMeasurements := board.Structures[description.Name].Measurements
	measurements := make([]MeasurementDTO, len(wantedMeasurements))
	for index, name := range wantedMeasurements {
		measurements[index] = board.Measurements[name]
	}

	return measurements
}

func getDescriptions(table domain.Table) map[Name]DescriptionDTO {
	descriptions := make(map[Name]DescriptionDTO, len(table.Rows))
	for _, row := range table.Rows {
		desc := newDescription(row)
		descriptions[desc.Name] = desc
	}

	return descriptions
}

func getMeasurements(table domain.Table) map[Name]MeasurementDTO {
	measurements := make(map[Name]MeasurementDTO, len(table.Rows))
	for _, row := range table.Rows {
		adapter := newMeasurement(row)
		measurements[adapter.Name] = adapter
	}

	return measurements
}

func getStructures(table domain.Table) map[Name]Structure {
	structures := make(map[Name]Structure)
	for _, column := range getColumns(table) {
		structure := newStructure(column)
		structures[structure.PacketName] = structure
	}

	return structures
}

func getColumns(table domain.Table) [][]string {
	columns := make([][]string, len(table.Rows[0]))
	for i := 0; i < len(table.Rows[0]); i++ {
		columns[i] = getColumn(i, table)
	}

	return columns
}

func getColumn(i int, table domain.Table) []string {
	column := make([]string, len(table.Rows))
	for j := 0; j < len(table.Rows); j++ {
		column[j] = table.Rows[j][i]
	}

	return column
}
