package domain

import (
	"time"
)

type PacketUpdate struct {
	ID            ID
	UpdatedValues map[string]any
	Timestamp     time.Time
}

func NewPacketUpdate(id ID, update map[string]any) PacketUpdate {
	return PacketUpdate{
		ID:            id,
		UpdatedValues: update,
		Timestamp:     time.Now(),
	}
}
