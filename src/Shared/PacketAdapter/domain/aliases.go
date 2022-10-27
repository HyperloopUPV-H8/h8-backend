package domain

type ID = uint16
type ValueType = string
type Numeric = float64
type EnumVariant = string
type Name = string

type PacketMeasurements = []MeasurementData

type MeasurementData struct {
	name      Name
	valueType ValueType
}
