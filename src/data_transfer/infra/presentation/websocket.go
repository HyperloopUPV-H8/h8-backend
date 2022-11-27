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
		buf := make(map[uint16]application.PacketJSON, 10)
		ticker := time.NewTicker(time.Millisecond * 10)
		for err == nil {
			payload := <-data
			buf[payload.ID] = payload
			select {
			case <-ticker.C:
				err = socket.WriteJSON(map[uint16]application.PacketJSON{payload.ID: payload})
			default:
			}

		}
	}()
}
