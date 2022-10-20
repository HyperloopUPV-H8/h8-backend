package adapters

import (
	excel "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/excelParser/domain"
)

type MeasurementAdapter struct {
	Name         string
	ValueType    string
	PodUnits     string
	DisplayUnits string
	SafeRange    string
	WarningRange string
}

func newMeasurementAdapter(row excel.Row) MeasurementAdapter {
	return MeasurementAdapter{
		Name:         row[0],
		ValueType:    row[1],
		PodUnits:     row[2],
		DisplayUnits: row[3],
		SafeRange:    row[4],
		WarningRange: row[5],
	}
}

func getMeasurementsWithSufix(sufix string, measurements []MeasurementAdapter) []MeasurementAdapter {
	values := make([]MeasurementAdapter, 0)
	for _, measurement := range measurements {
		newValue := measurement
		newValue.Name = newValue.Name + sufix
		values = append(values, newValue)
	}
	return values
}
