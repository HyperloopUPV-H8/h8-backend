package order_transfer

import (
	"encoding/json"
	"log"

	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
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

	return orderTransfer, channel
}

func initOrderTransfer() {
	orderChannel := make(chan vehicle_models.Order, ORDER_CHAN_BUFFER)
	orderTransfer = &OrderTransfer{orderChannel}
	channel = orderChannel
}

type OrderTransfer struct {
	channel chan<- vehicle_models.Order
}

func (orderTransfer *OrderTransfer) UpdateMessage(topic string, payload json.RawMessage, source string) {
	var order vehicle_models.Order
	if err := json.Unmarshal(payload, &order); err != nil {
		log.Printf("OrderTransfer: UpdateMessage: Unmarshal: %s\n", err)
		return
	}
	orderTransfer.channel <- order
}

func (orderTransfer *OrderTransfer) SetSendMessage(func(topic string, payload any, targets ...string) error) {
}

func (orderTransfer *OrderTransfer) HandlerName() string {
	return ORDER_TRASNFER_NAME
}
