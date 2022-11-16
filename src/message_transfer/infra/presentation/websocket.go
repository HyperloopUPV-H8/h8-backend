package presentation

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/server/infra/interfaces"
	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer/application"
)

func MessageRoutine(socket interfaces.WebSocket, messages <-chan application.MessageJSON) {
	go func() {
		defer socket.Close()
		var err error
		for err == nil {
			err = socket.WriteJSON(<-messages)
		}
	}()
}
