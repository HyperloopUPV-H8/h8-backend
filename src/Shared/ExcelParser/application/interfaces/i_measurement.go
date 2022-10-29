package interfaces

type Measurement interface {
	Name() string
	ValueType() string
	PodUnits() string
	DisplayUnits() string
	SafeRange() string
	WarningRange() string
}
