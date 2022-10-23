package measurement

import "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain/measurement/value"

type Measurement struct {
	Name   string
	Value  value.Value
	Units  Units
	Ranges Ranges
}
