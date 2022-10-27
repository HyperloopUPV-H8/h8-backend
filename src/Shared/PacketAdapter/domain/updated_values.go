package domain

import (
	"time"
)

type PacketUpdate struct {
	ID            uint16
	UpdatedValues map[string]any
	Timestamp     time.Time
}

func NewUpdatedValues(id ID, update map[string]any) PacketUpdate {
	return PacketUpdate{
		ID:            id,
		UpdatedValues: update,
		Timestamp:     time.Now(),
	}
}
