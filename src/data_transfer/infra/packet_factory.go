package infra

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/infra/interfaces"
)

type id = uint16

type PacketFactory struct {
	count     map[id]uint
	timestamp map[id]time.Time
}

func NewFactory() *PacketFactory {
	return &PacketFactory{
		count:     make(map[id]uint),
		timestamp: make(map[id]time.Time),
	}
}

func (factory *PacketFactory) NewPacket(data interfaces.Update) domain.Packet {
	cycleTime := data.Timestamp().Sub(factory.timestamp[data.ID()])
	factory.update(data.ID(), data.Timestamp())
	count := factory.count[data.ID()]
	return domain.NewPacket(data.ID(), count, cycleTime, data.HexValue(), data.Values())
}

func (factory *PacketFactory) update(id id, timestamp time.Time) {
	factory.count[id] += 1
	factory.timestamp[id] = timestamp
}
