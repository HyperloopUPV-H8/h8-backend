package domain

import (
	"io"

	value "github.com/HyperloopUPV-H8/Backend-H8/..."
)

type Packet struct {
	measurements []Measure
}

func (packet Packet) Decode(enums map[Name]Enum, bytes io.Reader) map[Name]value.Value {
	values := make(map[Name]value.Value, len(packet.measurements))
	for _, measure := range packet.measurements {
		values[measure.name] = measure.Decode(enums, bytes)
	}
	return values
}
