package vehicle

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/HyperloopUPV-H8/Backend-H8/pipe"
	"github.com/HyperloopUPV-H8/Backend-H8/sniffer"
	"github.com/HyperloopUPV-H8/Backend-H8/unit_converter"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/parsers"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/parsers/packet_parser"
	"github.com/rs/zerolog"
)

type Vehicle struct {
	sniffer sniffer.Sniffer
	pipes   map[string]*pipe.Pipe

	displayConverter unit_converter.UnitConverter
	podConverter     unit_converter.UnitConverter

	packetParser     packet_parser.PacketParser
	protectionParser parsers.ProtectionParser
	bitarrayParser   parsers.BitarrayParser

	dataChan chan packet.Raw

	idToBoard map[uint16]string

	onConnectionChange func(string, bool)

	trace zerolog.Logger
}

func (vehicle *Vehicle) Listen(updateChan chan<- models.PacketUpdate, protectionChan chan<- models.Protection) {
	vehicle.trace.Debug().Msg("vehicle listening")
	for raw := range vehicle.dataChan {
		payloadCopy := make([]byte, len(raw.Payload))
		copy(payloadCopy, raw.Payload)

		switch id := raw.Metadata.ID; {
		case vehicle.packetParser.Ids.Has(id):
			update, err := vehicle.packetParser.Decode(id, raw.Payload, raw.Metadata)

			if err != nil {
				vehicle.trace.Error().Err(err).Msg("error decoding packet")
				continue
			}

			convertedUpdate := vehicle.applyUnitConversion(update)

			updateChan <- convertedUpdate

		case vehicle.protectionParser.Ids.Has(id):
			protection, err := vehicle.protectionParser.Parse(id, raw.Payload)

			if err != nil {
				vehicle.trace.Error().Err(err).Msg("error decoding protection")
				continue
			}
			protectionChan <- protection
		default:
			fmt.Println("UNEXPECTED VALUE")
		}

		// switch payload := decoded.Payload.(type) {
		// case data.Payload:
		// 	vehicle.handleData(decoded.Metadata, payload, dataOutput)
		// case message.Payload:
		// 	vehicle.handleMessage(decoded.Metadata, payload, messageOutput)
		// case order.Payload:
		// 	vehicle.handleOrder(decoded.Metadata, payload, orderOutput)
		// default:
		// 	vehicle.trace.Error().Msg("unknown payload type")
		// }
	}
}

func (vehicle *Vehicle) applyUnitConversion(update models.PacketUpdate) models.PacketUpdate {
	newValues := make(map[string]packet.Value)

	for name, value := range update.Values {
		switch typedValue := value.(type) {
		case packet.Numeric:
			valueInSIUnits, podErr := vehicle.podConverter.Revert(name, float64(typedValue))

			if podErr != nil {
				//TODO: trace
			}

			valueInDisplayUnits, displayErr := vehicle.displayConverter.Convert(name, valueInSIUnits)

			if displayErr != nil {
				//TODO: trace
			}

			newValues[name] = packet.Numeric(valueInDisplayUnits)
		default:
			newValues[name] = typedValue
		}
	}

	return models.PacketUpdate{
		Metadata: update.Metadata,
		HexValue: update.HexValue,
		Values:   newValues,
	}
}

// func (vehicle *Vehicle) handleData(metadata packet.Metadata, payload data.Payload, output chan<- packet.Packet) {
// 	vehicle.trace.Trace().Uint16("id", metadata.ID).Msg("handle data")
// 	payload.Values = vehicle.podConverter.Revert(payload.Values)
// 	payload.Values = vehicle.displayConverter.Convert(payload.Values)

// 	select {
// 	case output <- packet.New(metadata, payload):
// 	default:
// 		vehicle.trace.Warn().Msg("data channel full")
// 	}
// }

// func (vehicle *Vehicle) handleMessage(metadata packet.Metadata, payload message.Payload, output chan<- packet.Packet) {
// 	vehicle.trace.Trace().Uint16("id", metadata.ID).Msg("handle message")

// 	select {
// 	case output <- packet.New(metadata, payload):
// 	default:
// 		vehicle.trace.Warn().Msg("message channel full")
// 	}
// }

// func (vehicle *Vehicle) handleOrder(metadata packet.Metadata, payload order.Payload, output chan<- packet.Packet) {
// 	vehicle.trace.Trace().Uint16("id", metadata.ID).Msg("handle order")

// 	payload.Values = vehicle.podConverter.Revert(payload.Values)
// 	payload.Values = vehicle.displayConverter.Convert(payload.Values)

// 	select {
// 	case output <- packet.New(metadata, payload):
// 	default:
// 		vehicle.trace.Warn().Msg("order channel full")
// 	}
// }

// func (vehicle *Vehicle) SendOrder(id uint16, vehicleOrder packet.Payload) (packet.Metadata, error) {
// 	vehicle.trace.Info().Uint16("id", id).Msg("send order")
// 	pipe, err := vehicle.getPipe(id)
// 	if err != nil {
// 		vehicle.trace.Error().Err(err).Msg("error getting pipe")
// 		return packet.Metadata{}, err
// 	}

// 	payload, ok := vehicleOrder.(order.Payload)
// 	if !ok {
// 		vehicle.trace.Error().Msg("payload is not order")
// 		return packet.Metadata{}, fmt.Errorf("payload is not order")
// 	}

// 	payload.Values = vehicle.displayConverter.Revert(payload.Values)
// 	payload.Values = vehicle.podConverter.Convert(payload.Values)

// 	data, err := vehicle.parser.Encode(id, payload)
// 	if err != nil {
// 		vehicle.trace.Error().Err(err).Msg("error encoding order")
// 		return packet.Metadata{}, err
// 	}

// 	_, err = common.WriteAll(pipe, data)

// 	return packet.Metadata{
// 		ID:        id,
// 		Timestamp: time.Now(),
// 		From:      pipe.Laddr(),
// 		To:        pipe.Raddr(),
// 		SeqNum:    0,
// 	}, err
// }

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
