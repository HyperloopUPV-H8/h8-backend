package domain

type MeasurementData struct {
	name      string
	valueType string
}

func NewMeasurement(name string, valueType string) MeasurementData {
	return MeasurementData{
		name:      name,
		valueType: valueType,
	}
}

func (measurement MeasurementData) Name() string {
	return measurement.name
}

func (measurement MeasurementData) ValueType() string {
	return measurement.valueType
}
