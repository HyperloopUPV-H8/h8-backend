package order_transfer

import (
	"encoding/json"

	"github.com/HyperloopUPV-H8/Backend-H8/common/observable"
	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const (
	ORDER_TRASNFER_NAME = "orderTransfer"
	ORDER_CHAN_BUFFER   = 100

	OrderTopic = "orders/enabled" // TODO: move to config
)

type OrderTransfer struct {
	stateOrders           map[string][]uint16
	stateOrdersObservable observable.ReplayObservable[map[string][]uint16]
	channel               chan<- vehicle_models.Order
	sendMessage           func(topic string, payload any, target ...string) error
	trace                 zerolog.Logger
}

func New() (OrderTransfer, <-chan vehicle_models.Order) {
	trace.Info().Msg("new order transfer")
	channel := make(chan vehicle_models.Order, ORDER_CHAN_BUFFER)
	stateOrders := make(map[string][]uint16)
	return OrderTransfer{
		channel:               channel,
		stateOrders:           stateOrders,
		stateOrdersObservable: observable.NewReplayObservable(stateOrders),
		trace:                 trace.With().Str("component", ORDER_TRASNFER_NAME).Logger(),
	}, channel
}

func (orderTransfer *OrderTransfer) UpdateMessage(topic string, payload json.RawMessage, source string) {
	orderTransfer.trace.Warn().Str("source", source).Str("topic", topic).Msg("got message")
	switch topic {
	case "order/send":
		orderTransfer.handleOrder(topic, payload, source)
	case "order/stateOrders":
		orderTransfer.handleSubscription(topic, payload, source)
	}
}

func (orderTransfer *OrderTransfer) handleSubscription(topic string, payload json.RawMessage, source string) {
	observable.HandleSubscribe[map[string][]uint16](&orderTransfer.stateOrdersObservable, source, payload, func(v map[string][]uint16, id string) error {
		return orderTransfer.sendMessage("order/stateOrders", v, id)
	})
}

func (orderTransfer *OrderTransfer) UpdateStateOrders(stateOrders vehicle_models.StateOrdersMessage) {
	orderTransfer.stateOrders[stateOrders.BoardId] = stateOrders.Orders
	orderTransfer.stateOrdersObservable.Next(orderTransfer.stateOrders)
}

func (orderTransfer *OrderTransfer) handleOrder(topic string, payload json.RawMessage, source string) {
	var order vehicle_models.Order
	if err := json.Unmarshal(payload, &order); err != nil {
		orderTransfer.trace.Error().Stack().Err(err).Msg("")
		return
	}
	orderTransfer.trace.Info().Str("source", source).Str("topic", topic).Uint16("id", order.ID).Msg("send order")
	orderTransfer.channel <- order
}

func (orderTransfer *OrderTransfer) SetSendMessage(sendMessage func(topic string, payload any, targets ...string) error) {
	orderTransfer.sendMessage = sendMessage
}

func (orderTransfer *OrderTransfer) HandlerName() string {
	return ORDER_TRASNFER_NAME
}
