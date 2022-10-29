package board

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/application/interfaces"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain"
)

type Measurement struct {
	name         string
	valueType    string
	podUnits     string
	displayUnits string
	safeRange    string
	warningRange string
}

func newMeasurement(row domain.Row) Measurement {
	return Measurement{
		name:         row[0],
		valueType:    row[1],
		podUnits:     row[2],
		displayUnits: row[3],
		safeRange:    row[4],
		warningRange: row[5],
	}
}

func (measurement Measurement) Name() string {
	return measurement.name
}

func (measurement Measurement) ValueType() string {
	return measurement.valueType
}

func (measurement Measurement) PodUnits() string {
	return measurement.podUnits
}

func (measurement Measurement) DisplayUnits() string {
	return measurement.displayUnits
}

func (measurement Measurement) SafeRange() string {
	return measurement.safeRange
}

func (measurement Measurement) WarningRange() string {
	return measurement.warningRange
}

func measurementsWithSufix(measurements []interfaces.Measurement, sufix string) []interfaces.Measurement {
	withSufix := make([]interfaces.Measurement, len(measurements))
	for i, measurement := range measurements {
		withSufix[i] = measurementWithName(measurement, measurement.Name()+sufix)
	}
	return withSufix
}

func measurementWithName(measurement interfaces.Measurement, name string) interfaces.Measurement {
	return Measurement{
		name:         name,
		valueType:    measurement.ValueType(),
		podUnits:     measurement.PodUnits(),
		displayUnits: measurement.DisplayUnits(),
		safeRange:    measurement.SafeRange(),
		warningRange: measurement.WarningRange(),
	}
}
