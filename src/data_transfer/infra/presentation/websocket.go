package presentation

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/server/infra/interfaces"
	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/application"
)

type id = uint16

func DataRoutine(socket interfaces.WebSocket, data <-chan application.PacketJSON) {
	go func() {
		var err error
		buf := make(map[uint16]application.PacketJSON)
		ticker := time.NewTicker(time.Millisecond * 10)
		for err == nil {
			select {
			case <-ticker.C:
				err = socket.WriteJSON(buf)
			case payload := <-data:
				buf[payload.ID] = payload
			}
		}
	}()
}
