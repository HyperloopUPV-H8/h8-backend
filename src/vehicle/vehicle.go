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
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/protection_parser"
	"github.com/rs/zerolog"
)

type Vehicle struct {
	sniffer sniffer.Sniffer
	pipes   map[string]*pipe.Pipe

	displayConverter unit_converter.UnitConverter
	podConverter     unit_converter.UnitConverter

	packetParser     packet_parser.PacketParser
	protectionParser protection_parser.ProtectionParser
	bitarrayParser   BitarrayParser
	orderIds         common.Set[uint16]

	dataChan chan packet.Packet

	idToBoard map[uint16]string

	onConnectionChange func(string, bool)

	trace zerolog.Logger
}

func (vehicle *Vehicle) Listen(updateChan chan<- models.PacketUpdate, transmittedOrderChan chan<- models.PacketUpdate, protectionChan chan<- models.ProtectionMessage, errorChan chan<- models.ErrorMessage) {
	vehicle.trace.Debug().Msg("vehicle listening")
	for packet := range vehicle.dataChan {
		payloadCopy := make([]byte, len(packet.Payload))
		copy(payloadCopy, packet.Payload)

		//TODO: add order decoding
		switch id := packet.Metadata.ID; {
		case vehicle.packetParser.Ids.Has(id):
			update, err := vehicle.getUpdate(packet)

			if err != nil {
				vehicle.trace.Error().Err(err).Msg("error decoding packet")
				continue
			}

			if vehicle.orderIds.Has(id) {
				transmittedOrderChan <- update
			} else {
				updateChan <- update
			}
		case vehicle.protectionParser.Ids.Has(id):
			message, err := vehicle.protectionParser.Parse(id, packet.Payload)

			if err != nil {
				vehicle.trace.Error().Err(err).Msg("error decoding protection")
				continue
			}
			protectionChan <- message
		default:
			vehicle.trace.Error().Uint16("id", packet.Metadata.ID).Msg("raw id not recognized")
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

	values := getOrderValues(order, vehicle.trace)
	convertedValues := vehicle.applyUnitConversion(values)

	buf := new(bytes.Buffer)

	idBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(idBuf, order.ID)

	err = vehicle.packetParser.Encode(order.ID, convertedValues, buf)
	if err != nil {
		vehicle.trace.Error().Err(err).Msg("error encoding order")
		return err
	}

	enableBuf := new(bytes.Buffer)
	vehicle.bitarrayParser.encodeBitarray(getOrderEnables(order), enableBuf)

	bufWithoutBitarray := append(idBuf, buf.Bytes()...)
	// fullBuf := append(bufWithoutBitarray, enableBuf.Bytes()...)

	_, err = common.WriteAll(pipe, bufWithoutBitarray)

	return err
}

func (vehicle *Vehicle) getUpdate(packet packet.Packet) (models.PacketUpdate, error) {
	update, err := vehicle.packetParser.Decode(packet.Metadata.ID, packet.Payload, packet.Metadata)

	if err != nil {
		return models.PacketUpdate{}, nil
	}

	convertedValues := vehicle.applyUnitConversion(update.Values)
	update.Values = convertedValues

	return update, nil
}

func (vehicle *Vehicle) applyUnitConversion(values map[string]packet.Value) map[string]packet.Value {
	newValues := make(map[string]packet.Value)

	for name, value := range values {
		switch typedValue := value.(type) {
		case packet.Numeric:
			newValues[name] = vehicle.applyNumericConversion(name, float64(typedValue))
		default:
			newValues[name] = typedValue
		}
	}

	return newValues
}

func (vehicle *Vehicle) applyNumericConversion(name string, value float64) packet.Numeric {
	valueInSIUnits, podErr := vehicle.podConverter.Revert(name, value)

	if podErr != nil {
		vehicle.trace.Error().Err(podErr).Msg("error reverting podUnits")
	}

	valueInDisplayUnits, displayErr := vehicle.displayConverter.Convert(name, valueInSIUnits)

	if displayErr != nil {
		vehicle.trace.Error().Err(displayErr).Msg("error converting to displayUnits")

	}

	return packet.Numeric(valueInDisplayUnits)
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

func getOrderValues(order models.Order, trace zerolog.Logger) map[string]packet.Value {
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
			trace.Error().Str("name", name).Type("type", field.Value).Msg("order field value not recognized")
		}
	}

	return values
}

func getOrderEnables(order models.Order) map[string]bool {
	enables := make(map[string]bool, 0)

	for name, field := range order.Fields {
		enables[name] = field.IsEnabled
	}

	return enables
}
