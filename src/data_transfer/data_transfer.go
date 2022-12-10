package data_transfer

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/models"
	"github.com/gorilla/websocket"
)

type DataTransfer struct {
	packetBuf map[uint16]models.PacketUpdate
	rate      time.Duration
}

func New(rate time.Duration) *DataTransfer {
	return &DataTransfer{
		packetBuf: make(map[uint16]models.PacketUpdate),
		rate:      rate,
	}
}

func (dataTransfer *DataTransfer) HandleConn(socket *websocket.Conn) {
	go func(socket *websocket.Conn) {
		defer socket.Close()
		ticker := time.NewTicker(dataTransfer.rate)

		for {
			<-ticker.C
			if err := socket.WriteJSON(dataTransfer.packetBuf); err != nil {
				break
			}
		}
	}(socket)
}

func (dataTransfer *DataTransfer) Update(update models.PacketUpdate) {
	dataTransfer.packetBuf[update.ID] = update
}
