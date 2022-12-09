package data_transfer

import (
	"fmt"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/models"
)

type PacketFactory struct {
	count     map[uint16]uint64
	timestamp map[uint16]uint64
}

func NewFactory() PacketFactory {
	return PacketFactory{
		count:     make(map[uint16]uint64),
		timestamp: make(map[uint16]uint64),
	}
}

func (factory PacketFactory) NewPacket(id uint16, raw []byte, values map[string]any) models.PacketUpdate {
	count, cycleTime := factory.getNext(id)
	return models.PacketUpdate{
		ID:        id,
		HexValue:  fmt.Sprintf("%x", raw),
		Values:    formatValues(values),
		Count:     count,
		CycleTime: cycleTime,
	}
}

func formatValues(values map[string]any) map[string]string {
	mapped := make(map[string]string, len(values))
	for name, val := range values {
		mapped[name] = fmt.Sprintf("%v", val)
	}
	return mapped
}

func (factory PacketFactory) getNext(id uint16) (count uint64, cycleTime uint64) {
	timestamp := uint64(time.Now().UnixMicro())
	cycleTime = timestamp - factory.timestamp[id]
	factory.timestamp[id] = timestamp
	factory.count[id] += 1
	count = factory.count[id]
	return
}
