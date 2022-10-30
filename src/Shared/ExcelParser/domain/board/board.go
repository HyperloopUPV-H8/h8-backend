package board

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain/document"
)

type Board struct {
	Name         string
	Descriptions map[Name]Description
	Measurements map[Name]Measurement
	Structures   map[Name]Structure
}

func New(sheet document.Sheet) Board {
	return Board{
		Name:         sheet.Name,
		Descriptions: getDescriptions(sheet.Tables["PacketDescription"]),
		Measurements: getMeasurements(sheet.Tables["ValueDescription"]),
		Structures:   getStructures(sheet.Tables["PacketStructure"]),
	}
}

func (board Board) GetPackets() []Packet {
	expandedPackets := make([]Packet, 0)
	for _, description := range board.Descriptions {
		measurements := board.getPacketMeasurements(description)
		packetDTOs := expandPacket(description, measurements)
		expandedPackets = append(expandedPackets, packetDTOs...)
	}
	return expandedPackets
}

func (board Board) getPacketMeasurements(description Description) []Measurement {
	wantedMeasurements := board.Structures[description.Name].Measurements
	measurements := make([]Measurement, len(wantedMeasurements))
	for index, name := range wantedMeasurements {
		measurements[index] = board.Measurements[name]
	}

	return measurements
}

func getDescriptions(table document.Table) map[Name]Description {
	descriptions := make(map[Name]Description, len(table.Rows))
	for _, row := range table.Rows {
		desc := newDescription(row)
		descriptions[desc.Name] = desc
	}

	return descriptions
}

func getMeasurements(table document.Table) map[Name]Measurement {
	measurements := make(map[Name]Measurement, len(table.Rows))
	for _, row := range table.Rows {
		adapter := newMeasurement(row)
		measurements[adapter.Name] = adapter
	}

	return measurements
}

func getStructures(table document.Table) map[Name]Structure {
	structures := make(map[Name]Structure)
	for _, column := range getColumns(table) {
		structure := newStructure(column)
		structures[structure.PacketName] = structure
	}

	return structures
}

func getColumns(table document.Table) [][]string {
	fmt.Println(table)
	columns := make([][]string, len(table.Rows[0]))
	for i := 0; i < len(table.Rows[0]); i++ {
		columns[i] = getColumn(i, table)
	}

	return columns
}

func getColumn(i int, table document.Table) []string {
	column := make([]string, len(table.Rows))
	for j := 0; j < len(table.Rows); j++ {
		column[j] = table.Rows[j][i]
	}

	return column
}
