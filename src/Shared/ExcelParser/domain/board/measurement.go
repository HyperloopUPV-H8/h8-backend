package board

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain/document"
)

type Measurement struct {
	Name         string
	ValueType    string
	PodUnits     string
	DisplayUnits string
	SafeRange    string
	WarningRange string
}

func newMeasurement(row document.Row) Measurement {
	return Measurement{
		Name:         row[0],
		ValueType:    row[1],
		PodUnits:     row[2],
		DisplayUnits: row[3],
		SafeRange:    row[4],
		WarningRange: row[5],
	}
}

func measurementsWithSufix(measurements []Measurement, sufix string) []Measurement {
	withSufix := make([]Measurement, len(measurements))
	for i, measurement := range measurements {
		withSufix[i] = measurementWithName(measurement, measurement.Name+sufix)
	}
	return withSufix
}

func measurementWithName(measurement Measurement, name string) Measurement {
	return Measurement{
		Name:         name,
		ValueType:    measurement.ValueType,
		PodUnits:     measurement.PodUnits,
		DisplayUnits: measurement.DisplayUnits,
		SafeRange:    measurement.SafeRange,
		WarningRange: measurement.WarningRange,
	}
}
