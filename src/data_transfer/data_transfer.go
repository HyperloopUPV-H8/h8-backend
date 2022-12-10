package data_transfer

import (
	"sync"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/models"
	"github.com/gorilla/websocket"
)

type DataTransfer struct {
	bufMx     sync.Mutex
	packetBuf map[uint16]models.PacketUpdate
	rate      time.Duration
}

func New(rate time.Duration) *DataTransfer {
	return &DataTransfer{
		bufMx:     sync.Mutex{},
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
			dataTransfer.bufMx.Lock()
			if err := socket.WriteJSON(dataTransfer.packetBuf); err != nil {
				dataTransfer.bufMx.Unlock()
				break
			}
			dataTransfer.bufMx.Unlock()
		}
	}(socket)
}

func (dataTransfer *DataTransfer) Update(update models.PacketUpdate) {
	dataTransfer.bufMx.Lock()
	dataTransfer.packetBuf[update.ID] = update
	dataTransfer.bufMx.Unlock()
}
