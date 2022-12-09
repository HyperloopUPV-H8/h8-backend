package data_transfer

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/models"
	"github.com/gorilla/websocket"
)

type DataTransfer struct {
	routines []chan<- models.PacketUpdate
	Rate     time.Duration
}

func (dataTransfer *DataTransfer) HandleConn(socket *websocket.Conn) {
	updates := make(chan models.PacketUpdate)
	dataTransfer.routines = append(dataTransfer.routines, updates)

	go func(socket *websocket.Conn, updates <-chan models.PacketUpdate, rate time.Duration) {
		defer socket.Close()
		buf := make(map[uint16]models.PacketUpdate)
		ticker := time.NewTicker(rate)
	loop:
		for {
			select {
			case payload := <-updates:
				buf[payload.ID] = payload
			case <-ticker.C:
				if err := socket.WriteJSON(buf); err != nil {
					break loop
				}
			}
		}
	}(socket, updates, dataTransfer.Rate)
}

func (dataTransfer *DataTransfer) Broadcast(update models.PacketUpdate) {
	for _, routine := range dataTransfer.routines {
		select {
		case routine <- update:
		default:
		}
	}
}
