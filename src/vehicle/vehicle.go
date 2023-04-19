package vehicle

import (
	"fmt"

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

func (vehicle *Vehicle) SendOrder(vehicleOrder models.Order) error {
	vehicle.trace.Info().Uint16("id", vehicleOrder.ID).Msg("send order")
	pipe, err := vehicle.getPipe(vehicleOrder.ID)
	if err != nil {
		return err
	}

	fields, enabled := unzipFields(vehicleOrder.Fields)
	fields = vehicle.displayConverter.Revert(fields)
	fields = vehicle.podConverter.Convert(fields)

	data, err := vehicle.parser.Encode(vehicleOrder.ID, order.Payload{Values: fields, Enabled: enabled})

	_, err = common.WriteAll(pipe, data)
	if err != nil {
		return err
	}

	return err
}

func (vehicle *Vehicle) getPipe(id uint16) (*pipe.Pipe, error) {
	board, ok := vehicle.idToBoard[id]
	if !ok {
		return nil, fmt.Errorf("board for id %d not found", id)
	}

	pipe, ok := vehicle.pipes[board]
	if !ok {
		return nil, fmt.Errorf("pipe for board %s not found", board)
	}

	return pipe, nil
}

func unzipFields(fields map[string]models.Field) (map[string]packet.Value, map[string]bool) {
	fieldsMap := make(map[string]packet.Value)
	enabledMap := make(map[string]bool)

	for name, field := range fields {
		fieldsMap[name] = field.Value
		enabledMap[name] = field.IsEnabled
	}

	return fieldsMap, enabledMap
}

func fieldsEnableToMap(fields map[string]models.Field) map[string]bool {
	enableMap := make(map[string]bool)

	for name, field := range fields {
		enableMap[name] = field.IsEnabled
	}

	return enableMap
}
