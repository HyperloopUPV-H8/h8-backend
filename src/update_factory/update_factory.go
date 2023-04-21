package update_factory

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"

	"github.com/HyperloopUPV-H8/Backend-H8/update_factory/models"
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

func (factory UpdateFactory) NewUpdate(update vehicle_models.PacketUpdate) models.Update {
	return models.Update{
		ID:        update.Metadata.ID,
		HexValue:  fmt.Sprintf("%x", update.HexValue),
		Values:    factory.getFields(update.Metadata.ID, update.Values),
		Count:     factory.getCount(update.Metadata.ID),
		CycleTime: factory.getCycleTime(update.Metadata.ID, uint64(update.Metadata.Timestamp.UnixNano())),
	}
}

func (factory UpdateFactory) getCount(id uint16) uint64 {
	if _, ok := factory.count[id]; !ok {
		factory.count[id] = 0
	}

	return factory.count[id]
}

func (factory UpdateFactory) getFields(id uint16, fields map[string]packet.Value) map[string]models.UpdateValue {
	updateFields := make(map[string]models.UpdateValue, len(fields))

	for name, value := range fields {
		switch value := value.(type) {
		case packet.Numeric:
			updateFields[name] = factory.getNumericField(id, name, float64(value))
		case packet.Boolean:
			updateFields[name] = models.BooleanValue(value)
		case packet.Enum:
			updateFields[name] = models.EnumValue(value)
		}
	}

	return updateFields
}

func (factory UpdateFactory) getNumericField(id uint16, name string, value float64) models.UpdateValue {
	avg := factory.getAverage(id, name)

	return models.NumericValue{Value: value, Average: avg.Add(value)}
}

func (factory UpdateFactory) getAverage(id uint16, name string) *common.MovingAverage[float64] {
	averages, ok := factory.fieldAvg[id]
	if !ok {
		averages = make(map[string]*common.MovingAverage[float64])
		factory.fieldAvg[id] = averages
	}

	average, ok := averages[name]
	if !ok {
		average = common.NewMovingAverage[float64](DEFAULT_ORDER)
		averages[name] = average
	}

	return average
}

func (factory UpdateFactory) getCycleTime(id uint16, timestamp uint64) uint64 {
	average, ok := factory.cycleTimeAvg[id]
	if !ok {
		average = common.NewMovingAverage[uint64](DEFAULT_ORDER)
		factory.cycleTimeAvg[id] = average
	}

	last, ok := factory.timestamp[id]
	if !ok {
		last = timestamp
		factory.timestamp[id] = last
	}

	cycleTime := timestamp - last
	factory.timestamp[id] = timestamp

	return average.Add(cycleTime)
}
