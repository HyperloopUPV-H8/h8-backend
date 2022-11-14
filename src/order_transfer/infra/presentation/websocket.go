package presentation

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/server/infra/interfaces"
	"github.com/HyperloopUPV-H8/Backend-H8/order_transfer/domain"
)

func OrderRoutine(socket interfaces.WebSocket, orders chan<- domain.Order) {
	go func() {
		var err error
		for err == nil {
			var ord domain.Order
			err = socket.ReadJSON(ord)
			orders <- ord
		}
	}()
}
