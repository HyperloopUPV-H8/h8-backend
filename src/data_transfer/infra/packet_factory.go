package infra

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/infra/dto"
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

func (factory *PacketFactory) NewPacket(data interfaces.Update) dto.Packet {
	return dto.NewPacket(factory.getCount(data.ID()), factory.getCycleTime(data.ID(), data.Timestamp()), data)
}

func (factory *PacketFactory) getCount(id id) uint {
	factory.count[id] += 1
	return factory.count[id]
}

func (factory *PacketFactory) getCycleTime(id id, timestamp time.Time) (cycleTime time.Duration) {
	cycleTime = timestamp.Sub(factory.timestamp[id])
	factory.timestamp[id] = timestamp
	return cycleTime
}
