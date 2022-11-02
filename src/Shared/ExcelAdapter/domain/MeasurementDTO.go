package domain

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelRetriever/domain"
)

type MeasurementDTO struct {
	Name         string
	ValueType    string
	PodUnits     string
	DisplayUnits string
	SafeRange    string
	WarningRange string
}

func newMeasurement(row domain.Row) MeasurementDTO {
	return MeasurementDTO{
		Name:         row[0],
		ValueType:    row[1],
		PodUnits:     row[2],
		DisplayUnits: row[3],
		SafeRange:    row[4],
		WarningRange: row[5],
	}
}

func measurementsWithSufix(measurements []MeasurementDTO, sufix string) []MeasurementDTO {
	withSufix := make([]MeasurementDTO, len(measurements))
	for i, measurement := range measurements {
		withSufix[i] = measurementWithName(measurement, measurement.Name+sufix)
	}
	return withSufix
}

func measurementWithName(measurement MeasurementDTO, name string) MeasurementDTO {
	return MeasurementDTO{
		Name:         name,
		ValueType:    measurement.ValueType,
		PodUnits:     measurement.PodUnits,
		DisplayUnits: measurement.DisplayUnits,
		SafeRange:    measurement.SafeRange,
		WarningRange: measurement.WarningRange,
	}
}
