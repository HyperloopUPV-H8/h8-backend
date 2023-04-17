package vehicle

import (
	"fmt"
	"log"
	"sync"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/message_parser"
	"github.com/HyperloopUPV-H8/Backend-H8/packet_parser"
	"github.com/HyperloopUPV-H8/Backend-H8/pipe"
	"github.com/HyperloopUPV-H8/Backend-H8/sniffer"
	"github.com/HyperloopUPV-H8/Backend-H8/unit_converter"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/internals"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/rs/zerolog"
)

type Vehicle struct {
	sniffer          sniffer.Sniffer
	parser           packet_parser.PacketParser
	messageParser    message_parser.MessageParser
	displayConverter unit_converter.UnitConverter
	podConverter     unit_converter.UnitConverter
	pipes            map[string]*pipe.Pipe

	packetFactory internals.UpdateFactory

	updateChan  chan []byte
	messageChan chan []byte

	idToBoard map[uint16]string

	stats              Stats
	statsMx            *sync.Mutex
	onConnectionChange func(string, bool)

	trace zerolog.Logger
}

func (vehicle *Vehicle) Listen(updateChan chan<- models.Update, messagesChan chan<- interface{}) {
	vehicle.trace.Info().Msg("start listening")
	go func() {
		for raw := range vehicle.updateChan {
			rawCopy := make([]byte, len(raw))
			copy(rawCopy, raw)

			id, fields := vehicle.parser.Decode(rawCopy)
			fields = vehicle.podConverter.Revert(fields)
			fields = vehicle.displayConverter.Convert(fields)

			update := vehicle.packetFactory.NewUpdate(id, rawCopy, fields)
			vehicle.statsMx.Lock()
			vehicle.stats.recv++
			vehicle.statsMx.Unlock()

			vehicle.trace.Trace().Msg("read")
			updateChan <- update
		}
	}()

	go func() {
		for raw := range vehicle.messageChan {
			msg, err := vehicle.messageParser.Parse(raw)
			if err != nil {
				vehicle.trace.Error().Stack().Err(err).Str("raw", fmt.Sprintf("%#v", string(raw))).Msg("parse message")
				continue
			}
			messagesChan <- msg
		}
	}()

}

func (vehicle *Vehicle) SendOrder(order models.Order) error {
	vehicle.trace.Info().Uint16("id", order.ID).Msg("send order")
	pipe, ok := vehicle.pipes[vehicle.idToBoard[order.ID]]
	if !ok {
		err := fmt.Errorf("%s pipe for %d not found", vehicle.idToBoard[order.ID], order.ID)
		vehicle.trace.Error().Stack().Err(err).Msg("")
		return err
	}

	valuesMap := fieldsValuesToMap(order.Fields)
	valuesMap = vehicle.displayConverter.Revert(valuesMap)
	valuesMap = vehicle.podConverter.Convert(valuesMap)

	valuesRaw := vehicle.parser.Encode(order.ID, valuesMap)
	bitArray := vehicle.parser.CreateBitArray(order.ID, fieldsEnableToMap(order.Fields))
	fullRaw := append(valuesRaw, bitArray...)
	log.Println("fullRaw", fullRaw)

	_, err := common.WriteAll(pipe, fullRaw)

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

func fieldsValuesToMap(fields map[string]models.Field) map[string]any {
	valuesMap := make(map[string]any)

	for name, field := range fields {
		valuesMap[name] = field.Value
	}

	return valuesMap
}

func fieldsEnableToMap(fields map[string]models.Field) map[string]bool {
	enableMap := make(map[string]bool)

	for name, field := range fields {
		enableMap[name] = field.IsEnabled
	}

	return enableMap
}
