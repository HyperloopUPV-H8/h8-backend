package domain

import (
	"time"
)

type id = uint16

type Packet struct {
	id        id
	count     uint
	cycleTime time.Duration
	hexValue  []byte
	values    map[string]any
}

func NewPacket(id id, count uint, cycleTime time.Duration, hexValue []byte, values map[string]any) Packet {
	return Packet{
		id:        id,
		count:     count,
		cycleTime: cycleTime,
		hexValue:  hexValue,
		values:    values,
	}
}

func (packet Packet) ID() id {
	return packet.id
}

func (packet Packet) Count() uint {
	return packet.count
}

func (packet Packet) CycleTime() time.Duration {
	return packet.cycleTime
}

func (packet Packet) HexValue() []byte {
	return packet.hexValue
}

func (packet Packet) Values() map[string]any {
	return packet.values
}
