package measurement

import (
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain/measurement/value"
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
)

type Measurement struct {
	Name   string
	Value  value.Value
	Ranges Ranges
}

func (m *Measurement) getDisplayString() string {
	return m.Value.ToDisplayUnitsString()
}

func NewMeasurements(rawMeasurements []excelAdapter.MeasurementDTO) map[string]Measurement {
	measurements := make(map[string]Measurement, len(rawMeasurements))
	for _, measurement := range rawMeasurements {
		measurements[measurement.Name] = Measurement{
			Name:   measurement.Name,
			Value:  value.NewDefault(measurement.ValueType, measurement.PodUnits, measurement.DisplayUnits),
			Ranges: NewRanges(measurement.SafeRange, measurement.WarningRange),
		}
	}
	return measurements
}
