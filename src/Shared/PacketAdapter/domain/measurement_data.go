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
