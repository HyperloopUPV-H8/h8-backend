package order_transfer

import (
	"encoding/json"
	"log"

	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	ws_models "github.com/HyperloopUPV-H8/Backend-H8/websocket_handle/models"
)

type OrderTransfer struct {
	orderChannel chan<- models.Order
	channel      chan ws_models.MessageTarget
}

func New(channel chan<- models.Order) (*OrderTransfer, chan ws_models.MessageTarget) {
	orderTransfer := &OrderTransfer{
		orderChannel: channel,
		channel:      make(chan ws_models.MessageTarget),
	}

	go orderTransfer.run()

	return orderTransfer, orderTransfer.channel
}

func (orderTransfer *OrderTransfer) run() {
	for msg := range orderTransfer.channel {
		log.Println(string(msg.Msg.Msg))
		var order models.Order
		err := json.Unmarshal(msg.Msg.Msg, &order)
		if err != nil {
			log.Printf("orderTransfer: run: %s\n", err)
			continue
		}
		orderTransfer.orderChannel <- order
	}
}
