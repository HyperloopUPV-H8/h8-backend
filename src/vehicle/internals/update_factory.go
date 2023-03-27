package internals

import (
	"fmt"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

type UpdateFactory struct {
	count     map[uint16]uint64
	timestamp map[uint16]uint64
}

func NewFactory() *UpdateFactory {
	return &UpdateFactory{
		count:     make(map[uint16]uint64),
		timestamp: make(map[uint16]uint64),
	}
}

func (factory UpdateFactory) NewUpdate(id uint16, hexValue []byte, fields map[string]any) models.Update {
	count, cycleTime := factory.getNext(id)
	return models.Update{
		ID:        id,
		HexValue:  fmt.Sprintf("%x", hexValue),
		Fields:    fields,
		Count:     count,
		CycleTime: cycleTime,
	}
}

func (factory UpdateFactory) getNext(id uint16) (count uint64, cycleTime uint64) {
	timestamp := uint64(time.Now().UnixMicro())
	cycleTime = timestamp - factory.timestamp[id]
	factory.timestamp[id] = timestamp
	factory.count[id] += 1
	count = factory.count[id]
	return
}
