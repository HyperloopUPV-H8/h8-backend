package internals

import (
	"fmt"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const DEFAULT_ORDER = 100

type UpdateFactory struct {
	count        map[uint16]uint64
	cycleTimeAvg map[uint16]*common.MovingAverage[uint64]
	timestamp    map[uint16]uint64
	fieldAvg     map[uint16]map[string]*common.MovingAverage[float64]
	trace        zerolog.Logger
}

func NewFactory() UpdateFactory {
	trace.Info().Msg("new update factory")
	return UpdateFactory{
		count:        make(map[uint16]uint64),
		cycleTimeAvg: make(map[uint16]*common.MovingAverage[uint64]),
		timestamp:    make(map[uint16]uint64),
		fieldAvg:     make(map[uint16]map[string]*common.MovingAverage[float64]),
		trace:        trace.With().Str("component", "updateFactory").Logger(),
	}
}

func (factory UpdateFactory) NewUpdate(id uint16, hexValue []byte, fields map[string]any) models.Update {
	count, cycleTime, averages := factory.getNext(id, fields)
	factory.trace.Trace().Uint16("id", id).Uint64("count", count).Uint64("cycleTime", cycleTime).Msg("new update")

	return models.Update{
		ID:        id,
		HexValue:  fmt.Sprintf("%x", hexValue),
		Fields:    fields,
		Averages:  averages,
		Count:     count,
		CycleTime: cycleTime,
	}
}

func (factory UpdateFactory) getNext(id uint16, fields map[string]any) (count uint64, cycleTime uint64, averages map[string]any) {
	timestamp := uint64(time.Now().UnixMicro())

	cycleTime = factory.getCycleTime(id, timestamp)

	factory.timestamp[id] = timestamp

	factory.count[id] += 1
	count = factory.count[id]

	averages = factory.getAverages(id, fields)

	return count, cycleTime, averages
}

func (factory UpdateFactory) getCycleTime(id uint16, timestamp uint64) uint64 {
	if _, ok := factory.cycleTimeAvg[id]; !ok {
		movAvg := common.NewMovingAverage[uint64](DEFAULT_ORDER)
		factory.cycleTimeAvg[id] = &movAvg
	}

	if _, ok := factory.timestamp[id]; !ok {
		factory.timestamp[id] = timestamp
	}

	cycleTimeAvg := factory.cycleTimeAvg[id]
	return cycleTimeAvg.Add(timestamp - factory.timestamp[id])
}

func (factory UpdateFactory) getAverages(id uint16, fields map[string]any) map[string]any {
	averages := make(map[string]any, len(fields))
	for key, value := range fields {
		averages[key] = factory.getAverage(id, key, value)
	}
	return averages
}

func (factory UpdateFactory) getAverage(id uint16, key string, value any) any {
	switch value := value.(type) {
	case float64:
		return factory.getFloatAverage(id, key, value)
	default:
		return value
	}
}

func (factory UpdateFactory) getFloatAverage(id uint16, key string, value float64) float64 {
	if _, ok := factory.fieldAvg[id]; !ok {
		factory.fieldAvg[id] = make(map[string]*common.MovingAverage[float64])
	}

	if _, ok := factory.fieldAvg[id][key]; !ok {
		movAvg := common.NewMovingAverage[float64](DEFAULT_ORDER)
		factory.fieldAvg[id][key] = &movAvg
	}

	fieldAvg := factory.fieldAvg[id][key]
	return fieldAvg.Add(value)
}
