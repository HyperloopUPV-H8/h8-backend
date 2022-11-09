package dto

import (
	"time"
)

type id = uint16

type PacketUpdate struct {
	id        uint16
	hexValue  []byte
	values    map[string]any
	timestamp time.Time
}

func NewPacketUpdate(id uint16, update map[string]any, hexValues []byte) PacketUpdate {
	return PacketUpdate{
		id:        id,
		hexValue:  hexValues,
		values:    update,
		timestamp: time.Now(),
	}
}

func (update PacketUpdate) ID() id {
	return update.id
}

func (update PacketUpdate) HexValue() []byte {
	return update.hexValue
}

func (update PacketUpdate) Values() map[string]any {
	return update.values
}

func (update PacketUpdate) Timestamp() time.Time {
	return update.timestamp
}
