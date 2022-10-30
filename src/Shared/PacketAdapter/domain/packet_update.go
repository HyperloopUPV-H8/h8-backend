package domain

import (
	"time"
)

type PacketUpdate struct {
	id            ID
	updatedValues map[string]any
	timestamp     time.Time
}

func NewPacketUpdate(id ID, update map[string]any) PacketUpdate {
	return PacketUpdate{
		id:            id,
		updatedValues: update,
		timestamp:     time.Now(),
	}
}

func (update PacketUpdate) ID() ID {
	return update.id
}

func (update PacketUpdate) Timestamp() time.Time {
	return update.timestamp
}

func (update PacketUpdate) Values() map[string]any {
	return update.updatedValues
}
