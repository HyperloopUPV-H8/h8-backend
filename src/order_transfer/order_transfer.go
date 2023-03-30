package order_transfer

import (
	"encoding/json"

	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const (
	ORDER_TRASNFER_NAME = "orderTransfer"
	ORDER_CHAN_BUFFER   = 100
)

var (
	orderTransfer *OrderTransfer
	channel       <-chan vehicle_models.Order
)

func Get() (*OrderTransfer, <-chan vehicle_models.Order) {
	if orderTransfer == nil {
		initOrderTransfer()
	}

	trace.Debug().Msg("get order transfer")
	return orderTransfer, channel
}

func initOrderTransfer() {
	trace.Info().Msg("init order transfer")
	orderChannel := make(chan vehicle_models.Order, ORDER_CHAN_BUFFER)
	orderTransfer = &OrderTransfer{
		channel: orderChannel,
		trace:   trace.With().Str("component", ORDER_TRASNFER_NAME).Logger(),
	}
	channel = orderChannel
}

type OrderTransfer struct {
	channel chan<- vehicle_models.Order
	trace   zerolog.Logger
}

func (orderTransfer *OrderTransfer) UpdateMessage(topic string, payload json.RawMessage, source string) {
	orderTransfer.trace.Warn().Str("source", source).Str("topic", topic).Msg("got message")
	var order vehicle_models.Order
	if err := json.Unmarshal(payload, &order); err != nil {
		orderTransfer.trace.Error().Stack().Err(err).Msg("")
		return
	}
	orderTransfer.trace.Info().Str("source", source).Str("topic", topic).Uint16("id", order.ID).Msg("send order")
	orderTransfer.channel <- order
}

func (orderTransfer *OrderTransfer) SetSendMessage(func(topic string, payload any, targets ...string) error) {
	orderTransfer.trace.Debug().Msg("set send message")
}

func (orderTransfer *OrderTransfer) HandlerName() string {
	return ORDER_TRASNFER_NAME
}
