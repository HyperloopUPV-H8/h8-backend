package vehicle

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/HyperloopUPV-H8/Backend-H8/pipe"
	"github.com/HyperloopUPV-H8/Backend-H8/sniffer"
	"github.com/HyperloopUPV-H8/Backend-H8/unit_converter"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/packet_parser"
	"github.com/rs/zerolog"
)

type Vehicle struct {
	sniffer sniffer.Sniffer
	pipes   map[string]*pipe.Pipe

	displayConverter unit_converter.UnitConverter
	podConverter     unit_converter.UnitConverter

	packetParser     packet_parser.PacketParser
	protectionParser ProtectionParser
	bitarrayParser   BitarrayParser

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

			convertedValues := vehicle.applyUnitConversion(update.Values)
			update.Values = convertedValues

			updateChan <- update

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

	}
}

func (vehicle *Vehicle) SendOrder(order models.Order) error {
	vehicle.trace.Info().Uint16("id", order.ID).Msg("send order")
	pipe, err := vehicle.getPipe(order.ID)

	if err != nil {
		vehicle.trace.Error().Err(err).Msg("error getting pipe")
		return err
	}

	values := getOrderValues(order)
	convertedValues := vehicle.applyUnitConversion(values)

	buf := new(bytes.Buffer)

	idBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(idBuf, order.ID)

	err = vehicle.packetParser.Encode(order.ID, convertedValues, buf)
	if err != nil {
		vehicle.trace.Error().Err(err).Msg("error encoding order")
		return err
	}

	fullBuf := append(idBuf, buf.Bytes()...)

	_, err = common.WriteAll(pipe, fullBuf)

	return err
}

func (vehicle *Vehicle) applyUnitConversion(values map[string]packet.Value) map[string]packet.Value {
	newValues := make(map[string]packet.Value)

	for name, value := range values {
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

	return newValues
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

func getOrderValues(order models.Order) map[string]packet.Value {
	values := make(map[string]packet.Value)

	for name, field := range order.Fields {
		switch value := field.Value.(type) {
		case float64:
			values[name] = packet.Numeric(value)
		case bool:
			values[name] = packet.Boolean(value)
		case string:
			values[name] = packet.Enum(value)
		default:
			//TODO: trace
		}
	}

	return values
}

func getOrderEnables(order models.Order) []bool {
	enables := make([]bool, 0)

	for _, field := range order.Fields {
		enables = append(enables, field.IsEnabled)
	}

	return enables
}
