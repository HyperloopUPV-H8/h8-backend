package order_transfer

import (
	"github.com/HyperloopUPV-H8/Backend-H8/order_transfer/models"
	"github.com/gorilla/websocket"
)

type OrderTransfer struct {
	OrderChannel chan<- models.Order
}

func (orderTransfer *OrderTransfer) HandleConn(socket *websocket.Conn) {
	go func(socket *websocket.Conn, orderChannel chan<- models.Order) {
		defer socket.Close()
		for {
			var payload models.Order
			if err := socket.ReadJSON(payload); err != nil {
				return
			}
			orderChannel <- payload
		}
	}(socket, orderTransfer.OrderChannel)
}
