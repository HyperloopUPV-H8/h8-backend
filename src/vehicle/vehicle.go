package vehicle

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/packet_parser"
	"github.com/HyperloopUPV-H8/Backend-H8/pipe"
	"github.com/HyperloopUPV-H8/Backend-H8/sniffer"
	"github.com/HyperloopUPV-H8/Backend-H8/unit_converter"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/internals"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/rs/zerolog"
)

type Vehicle struct {
	sniffer          *sniffer.Sniffer
	parser           *packet_parser.PacketParser
	displayConverter *unit_converter.UnitConverter
	podConverter     *unit_converter.UnitConverter
	pipes            map[string]*pipe.Pipe

	packetFactory *internals.UpdateFactory

	readChan chan []byte

	idToPipe map[uint16]string

	stats *Stats

	onConnectionChange func(string, bool)

	trace zerolog.Logger
}

func (vehicle *Vehicle) SendOrder(order models.Order) error {
	vehicle.trace.Info().Uint16("id", order.ID).Msg("send order")
	pipe, ok := vehicle.pipes[vehicle.idToPipe[order.ID]]
	if !ok {
		err := fmt.Errorf("%s pipe for %d not found", vehicle.idToPipe[order.ID], order.ID)
		vehicle.trace.Error().Err(err).Msg("")
		return err
	}

	fields := order.Fields
	fields = vehicle.displayConverter.Revert(fields)
	fields = vehicle.podConverter.Revert(fields)
	raw := vehicle.parser.Encode(order.ID, fields)

	_, err := common.WriteAll(pipe, raw)

	if err == nil {
		vehicle.stats.sent++
	} else {
		vehicle.trace.Error().Err(err).Msg("")
		vehicle.stats.sentFail++
	}

	return err
}

func (vehicle *Vehicle) Listen(output chan<- models.Update) {
	vehicle.trace.Info().Msg("start listening")
	for raw := range vehicle.readChan {
		rawCopy := make([]byte, len(raw))
		copy(rawCopy, raw)

		id, fields := vehicle.parser.Decode(rawCopy)
		fields = vehicle.podConverter.Convert(fields)
		fields = vehicle.displayConverter.Convert(fields)

		update := vehicle.packetFactory.NewUpdate(id, rawCopy, fields)

		vehicle.stats.recv++

		vehicle.trace.Trace().Msg("read")
		output <- update
	}
}

func (vehicle *Vehicle) OnConnectionChange(callback func(string, bool)) {
	vehicle.trace.Debug().Msg("set on connection change")
	vehicle.onConnectionChange = callback
}
