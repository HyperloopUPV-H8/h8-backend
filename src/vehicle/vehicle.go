package vehicle

import (
	"fmt"
	"sync"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/data"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/message"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/order"
	"github.com/HyperloopUPV-H8/Backend-H8/pipe"
	"github.com/HyperloopUPV-H8/Backend-H8/sniffer"
	"github.com/HyperloopUPV-H8/Backend-H8/unit_converter"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/rs/zerolog"
)

type Vehicle struct {
	sniffer sniffer.Sniffer
	pipes   map[string]*pipe.Pipe

	parser           *packet.Parser
	displayConverter unit_converter.UnitConverter
	podConverter     unit_converter.UnitConverter

	dataChan chan packet.Raw

	idToBoard map[uint16]string

	stats              Stats
	statsMx            *sync.Mutex
	onConnectionChange func(string, bool)

	trace zerolog.Logger
}

func (vehicle *Vehicle) Listen(dataOutput, messageOutput, orderOutput chan<- packet.Packet) {
	for raw := range vehicle.dataChan {
		payloadCopy := make([]byte, len(raw.Payload))
		copy(payloadCopy, raw.Payload)

		decoded, err := vehicle.parser.Decode(raw.Metadata, payloadCopy)
		if err != nil {
			// FIXME: handle error
			panic("error decoding packet")
		}

		switch payload := decoded.Payload.(type) {
		case data.Payload:
			vehicle.handleData(decoded.Metadata, payload, dataOutput)
		case message.Payload:
			vehicle.handleMessage(decoded.Metadata, payload, messageOutput)
		case order.Payload:
			vehicle.handleOrder(decoded.Metadata, payload, orderOutput)
		}
	}
}

func (vehicle *Vehicle) handleData(metadata packet.Metadata, payload data.Payload, output chan<- packet.Packet) {
	payload.Values = vehicle.podConverter.Revert(payload.Values)
	payload.Values = vehicle.displayConverter.Convert(payload.Values)
	output <- packet.New(metadata, payload)
}

func (vehicle *Vehicle) handleMessage(metadata packet.Metadata, payload message.Payload, output chan<- packet.Packet) {
	output <- packet.New(metadata, payload)
}

func (vehicle *Vehicle) handleOrder(metadata packet.Metadata, payload order.Payload, output chan<- packet.Packet) {
	payload.Values = vehicle.podConverter.Revert(payload.Values)
	payload.Values = vehicle.displayConverter.Convert(payload.Values)
	output <- packet.New(metadata, payload)
}

func (vehicle *Vehicle) SendOrder(order models.Order) error {
	vehicle.trace.Info().Uint16("id", order.ID).Msg("send order")
	pipe, ok := vehicle.pipes[vehicle.idToBoard[order.ID]]
	if !ok {
		err := fmt.Errorf("%s pipe for %d not found", vehicle.idToBoard[order.ID], order.ID)
		vehicle.trace.Error().Stack().Err(err).Msg("")
		return err
	}

	fields := order.Fields
	fields = vehicle.displayConverter.Convert(fields)
	fields = vehicle.podConverter.Revert(fields)
	raw := vehicle.parser.Encode(order.ID, fields)

	_, err := common.WriteAll(pipe, raw)

	vehicle.statsMx.Lock()
	defer vehicle.statsMx.Unlock()

	if err == nil {
		vehicle.stats.sent++
	} else {
		vehicle.trace.Error().Stack().Err(err).Msg("")
		vehicle.stats.sentFail++
	}

	return err
}
