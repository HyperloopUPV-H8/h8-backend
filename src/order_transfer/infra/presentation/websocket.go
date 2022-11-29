package presentation

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/server/infra/interfaces"
	"github.com/HyperloopUPV-H8/Backend-H8/order_transfer/application"
)

func OrderRoutine(socket interfaces.WebSocket, orders chan<- application.OrderJSON) {
	go func() {
		defer socket.Close()
		var err error
		for err == nil {
			var ord application.OrderJSON
			err = socket.ReadJSON(ord)
			if err == nil {
				orders <- ord
			}
		}
	}()
}
