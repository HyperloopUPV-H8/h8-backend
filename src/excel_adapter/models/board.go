package models

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/internals/models"
	"golang.org/x/exp/slices"
)

const (
	PACKET_TABLE_NAME          = "Packets"
	MEASUREMENT_TABLE_NAME     = "Measurements"
	STRUCTURES_TABLE_NAME      = "Structures"
	CONTROL_SECTION_TABLE_NAME = "ControlSections"
)

type Board struct {
	Name            string
	IP              string
	Descriptions    map[string]Description
	Measurements    map[string]Value
	Structures      map[string]Structure
	ControlSections map[string]ControlSection
}

func NewBoard(name string, ip string, sheet models.Sheet) Board {
	return Board{
		Name:            name,
		IP:              ip,
		Descriptions:    getDescriptions(sheet.Tables[PACKET_TABLE_NAME]),
		Measurements:    getMeasurements(sheet.Tables[MEASUREMENT_TABLE_NAME]),
		Structures:      getStructures(sheet.Tables[STRUCTURES_TABLE_NAME]),
		ControlSections: getControlSections(sheet.Tables[CONTROL_SECTION_TABLE_NAME]),
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

func (board Board) getPacketMeasurements(description Description) []Value {
	wantedMeasurements := board.Structures[description.Name].Measurements
	measurements := make([]Value, len(wantedMeasurements))
	for index, name := range wantedMeasurements {
		measurements[index] = board.Measurements[name]
	}

	return measurements
}

func getDescriptions(table models.Table) map[string]Description {
	descriptions := make(map[string]Description, len(table.Rows))
	for _, row := range table.Rows {
		desc := newDescription(row)
		descriptions[desc.Name] = desc
	}

	return descriptions
}

func getMeasurements(table models.Table) map[string]Value {
	measurements := make(map[string]Value, len(table.Rows))
	for _, row := range table.Rows {
		adapter := newValue(row)
		measurements[adapter.Name] = adapter
	}

	return measurements
}

func getStructures(table models.Table) map[string]Structure {
	structures := make(map[string]Structure)
	for _, column := range getColumns(table) {
		structure := newStructure(column)
		structures[structure.PacketName] = structure
	}

	return structures
}

func getControlSections(table models.Table) map[string]ControlSection {
	controlSections := make(map[string]ControlSection)
	for colIndex := 0; colIndex < len(table.Rows[0]); colIndex += 2 {
		controlSection := make(ControlSection)
		sectionTitle := table.Rows[0][colIndex]
		for rowIndex := 1; rowIndex < len(table.Rows); rowIndex++ {
			if table.Rows[rowIndex][colIndex] != "" {
				controlSection[table.Rows[rowIndex][colIndex]] = table.Rows[rowIndex][colIndex+1]
			} else {
				break
			}
		}
		controlSections[sectionTitle] = controlSection
	}
	return controlSections
}

func getColumns(table models.Table) [][]string {
	columns := make([][]string, len(table.Rows[0]))
	for i := 0; i < len(table.Rows[0]); i++ {
		columns[i] = getColumn(i, table)
	}

	return columns
}

func getColumn(i int, table models.Table) []string {
	column := make([]string, len(table.Rows))
	for j := 0; j < len(table.Rows); j++ {
		column[j] = table.Rows[j][i]
	}

	return column
}

func (board *Board) FindContainingPacket(valueName string) string {
	for _, structure := range board.Structures {
		if slices.IndexFunc(structure.Measurements, func(name string) bool {
			return name == valueName
		}) != -1 {
			return structure.PacketName
		}
	}

	panic(fmt.Sprintf("valueName %s doesn't exist", valueName))
}
