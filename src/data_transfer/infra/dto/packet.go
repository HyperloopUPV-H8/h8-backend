package dto

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/infra/interfaces"
)

type id = uint16

type Packet struct {
	id           id
	count        uint
	cycleTime    time.Duration
	hexValue     []byte
	measurements map[string]any
}

func NewPacket(count uint, cycleTime time.Duration, update interfaces.Update) Packet {
	return Packet{
		id:           update.ID(),
		count:        count,
		cycleTime:    cycleTime,
		hexValue:     update.HexValue(),
		measurements: update.Measurements(),
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

func (packet Packet) Measurements() map[string]any {
	return packet.measurements
}
