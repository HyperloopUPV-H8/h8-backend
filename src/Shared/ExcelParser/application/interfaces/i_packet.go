package interfaces

type Packet interface {
	Description() Description
	Measurements() []Measurement
}
