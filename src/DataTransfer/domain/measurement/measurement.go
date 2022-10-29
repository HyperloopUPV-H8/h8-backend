package measurement

import (
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain/measurement/value"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/application/interfaces"
)

type Measurement struct {
	Name   string
	Value  value.Value
	Ranges Ranges
}

func (m *Measurement) getDisplayString() string {
	return m.Value.ToDisplayString()
}

func NewMeasurements(rawMeasurements []interfaces.Measurement) map[string]Measurement {
	measurements := make(map[string]Measurement, len(rawMeasurements))
	for _, measurement := range rawMeasurements {
		measurements[measurement.Name()] = Measurement{
			Name:   measurement.Name(),
			Value:  value.NewDefault(measurement.ValueType(), measurement.PodUnits(), measurement.DisplayUnits()),
			Ranges: NewRanges(measurement.SafeRange(), measurement.WarningRange()),
		}
	}
	return measurements
}
