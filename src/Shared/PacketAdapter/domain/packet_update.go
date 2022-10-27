package domain

import (
	"time"

	value "github.com/HyperloopUPV-H8/Backend-H8/..."
)

type PacketUpdate struct {
	id        uint16
	measures  map[string]value.Value
	timestamp time.Time
}

func NewPacketUpdate(id ID, enums []Enum, measures map[string]any) PacketUpdate {
	return PacketUpdate{
		id:        uint16(id),
		measures:  measures,
		timestamp: time.Now(),
	}
}
