package domain

type MeasurementData struct {
	Name      string
	ValueType string
}

func NewMeasurement(name string, valueType string) MeasurementData {
	return MeasurementData{
		Name:      name,
		ValueType: valueType,
	}
}
