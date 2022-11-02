package domain

import (
	"time"
)

type PacketUpdate struct {
	ID            uint16
	HexValue      []byte
	UpdatedValues map[string]any
	Timestamp     time.Time
}

func NewPacketUpdate(id uint16, update map[string]any, hexValues []byte) PacketUpdate {
	return PacketUpdate{
		ID:            id,
		HexValue:      hexValues,
		UpdatedValues: update,
		Timestamp:     time.Now(),
	}
}
