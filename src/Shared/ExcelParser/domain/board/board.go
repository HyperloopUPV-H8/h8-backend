package board

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/application/interfaces"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain"
)

type Board struct {
	name         string
	descriptions map[Name]interfaces.Description
	measurements map[Name]interfaces.Measurement
	structures   map[Name]interfaces.Structure
}

func (board Board) Name() string {
	return board.name
}

func (board Board) Descriptions() map[string]interfaces.Description {
	return board.descriptions
}

func (board Board) Measurements() map[string]interfaces.Measurement {
	return board.measurements
}

func (board Board) Structure() map[string]interfaces.Structure {
	return board.structures
}

func NewBoard(sheet domain.Sheet) interfaces.Board {
	return Board{
		name:         sheet.Name,
		descriptions: getDescriptions(sheet.Tables["Packet Description"]),
		measurements: getMeasurements(sheet.Tables["Value Description"]),
		structures:   getStructures(sheet.Tables["Packet Structure"]),
	}
}

func (b Board) GetPackets() []interfaces.Packet {
	expandedPackets := make([]interfaces.Packet, 0)
	for _, description := range b.descriptions {
		measurements := b.getPacketMeasurements(description)
		packetDTOs := expandPacket(description, measurements)
		expandedPackets = append(expandedPackets, packetDTOs...)
	}
	return expandedPackets
}

func (b Board) getPacketMeasurements(description interfaces.Description) []interfaces.Measurement {
	wantedMeasurements := b.structures[description.Name()].Measurements()
	measurements := make([]interfaces.Measurement, len(wantedMeasurements))
	for index, name := range wantedMeasurements {
		measurements[index] = b.measurements[name]
	}

	return measurements
}

func getDescriptions(table domain.Table) map[Name]interfaces.Description {
	descriptions := make(map[Name]interfaces.Description, len(table.Rows))
	for _, row := range table.Rows {
		adapter := newDescription(row)
		descriptions[adapter.Name()] = adapter
	}

	return descriptions
}

func getMeasurements(table domain.Table) map[Name]interfaces.Measurement {
	measurements := make(map[Name]interfaces.Measurement, len(table.Rows))
	for _, row := range table.Rows {
		adapter := newMeasurement(row)
		measurements[adapter.Name()] = adapter
	}

	return measurements
}

func getStructures(table domain.Table) map[Name]interfaces.Structure {
	structures := make(map[Name]interfaces.Structure)
	for _, column := range getColumns(table) {
		structure := newStructure(column)
		structures[structure.PacketName()] = structure
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
