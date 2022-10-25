package dto

import (
	excel "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/excelRetreiver"
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain/measurement"
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain/measurement/value"
)

type MeasurementDTO struct {
	Name         string
	ValueType    string
	PodUnits     string
	DisplayUnits string
	SafeRange    string
	WarningRange string
}

func (m *MeasurementDTO) toMeasurement() measurement.Measurement {
	return measurement.Measurement{
		Name:   m.Name,
		Value:  value.NewDefault(m.ValueType, m.PodUnits, m.DisplayUnits),
		Ranges: measurement.NewRanges(m.SafeRange, m.WarningRange),
	}
}

func newMeasurementDTO(row excel.Row) MeasurementDTO {
	return MeasurementDTO{
		Name:         row[0],
		ValueType:    row[1],
		PodUnits:     row[2],
		DisplayUnits: row[3],
		SafeRange:    row[4],
		WarningRange: row[5],
	}
}

func getMeasurementsWithSufix(sufix string, measurements []MeasurementDTO) []MeasurementDTO {
	mArr := make([]MeasurementDTO, len(measurements))
	for index, measurement := range measurements {
		measurement.Name += sufix
		mArr[index] = measurement
	}
	return mArr
}
