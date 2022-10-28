package measurement

import "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain/measurement/value"

type Measurement struct {
	Name   string
	Value  value.Value
	Ranges Ranges
}

func (m *Measurement) getDisplayString() string {
	return m.Value.ToDisplayString()
}
