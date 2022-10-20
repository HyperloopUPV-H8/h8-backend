package adapters

import (
	excel "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/excelParser/domain"
)

type BoardAdapter struct {
	descriptions map[string]DescriptionAdapter
	measurements map[string]MeasurementAdapter
	structures   map[string]StructureAdapter
}

func NewBoardAdapter(tables map[string]excel.Table) BoardAdapter {
	adapter := BoardAdapter{}
	adapter.addMeasurements(tables["ValueDescription"])
	adapter.addDescriptions(tables["PacketDescription"])
	adapter.addStructures(tables["PacketStructure"])
	return adapter
}

func (boardAdapter *BoardAdapter) GetExpandedPacketAdapters() []PacketAdapter {
	expandedPackets := make([]PacketAdapter, 0)
	for _, description := range boardAdapter.descriptions {
		measurements := boardAdapter.getPacketMeasurements(description)
		packets := expandPacket(description, measurements)
		expandedPackets = append(expandedPackets, packets...)
	}
	return expandedPackets
}

func (boardAdapter *BoardAdapter) getPacketMeasurements(description DescriptionAdapter) []MeasurementAdapter {
	measurements := make([]MeasurementAdapter, 0)

	for _, name := range boardAdapter.structures[description.Name].measurements {
		measurements = append(measurements, boardAdapter.measurements[name])
	}

	return measurements
}

func (boardAdapter *BoardAdapter) addDescriptions(table excel.Table) {
	descriptions := make(map[string]DescriptionAdapter)
	for _, row := range table.Rows {
		adapter := newDescriptionAdapter(row)
		descriptions[adapter.Name] = adapter
	}

	boardAdapter.descriptions = descriptions
}

func (boardAdapter *BoardAdapter) addMeasurements(table excel.Table) {
	measurements := make(map[string]MeasurementAdapter)
	for _, row := range table.Rows {
		adapter := newMeasurementAdapter(row)
		measurements[adapter.Name] = adapter
	}

	boardAdapter.measurements = measurements
}

func (boardAdapter *BoardAdapter) addStructures(table excel.Table) {
	columns := make([][]string, 0)
	for i := 0; i < len(table.Rows[0]); i++ {
		column := make([]string, 0)
		for j := 0; j < len(table.Rows); j++ {
			column = append(column, table.Rows[j][i])
		}
		columns = append(columns, column)
	}

	for _, column := range columns {
		structure := newStructure(column)
		boardAdapter.structures[structure.packetName] = structure
	}
}
