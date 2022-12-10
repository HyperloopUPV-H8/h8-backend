package data_transfer

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/models"
	"github.com/gorilla/websocket"
	"github.com/kjk/betterguid"
)

type DataTransfer struct {
	bufMx     sync.Mutex
	packetBuf map[uint16]models.PacketUpdate
	ticker    *time.Ticker
	sockets   map[string]*websocket.Conn
}

func New(rate time.Duration) *DataTransfer {
	dataTransfer := &DataTransfer{
		bufMx:     sync.Mutex{},
		packetBuf: make(map[uint16]models.PacketUpdate),
		ticker:    time.NewTicker(rate),
		sockets:   make(map[string]*websocket.Conn),
	}

	go dataTransfer.run()

	return dataTransfer
}

func (dataTransfer *DataTransfer) run() {
	for {
		<-dataTransfer.ticker.C
		data := dataTransfer.getJSON()
		for id, socket := range dataTransfer.sockets {
			if err := socket.WriteMessage(websocket.TextMessage, data); err != nil {
				socket.Close()
				delete(dataTransfer.sockets, id)
			}
		}
	}
}

func (dataTransfer *DataTransfer) Close() {
	for _, socket := range dataTransfer.sockets {
		socket.Close()
	}
}

func (dataTransfer *DataTransfer) getJSON() []byte {
	dataTransfer.bufMx.Lock()
	defer dataTransfer.bufMx.Unlock()
	data, err := json.Marshal(dataTransfer.packetBuf)
	if err != nil {
		log.Fatalf("data transfer: getJSON: %s\n", err)
	}
	return data
}

func (dataTransfer *DataTransfer) HandleConn(socket *websocket.Conn) {
	dataTransfer.sockets[betterguid.New()] = socket
}

func (dataTransfer *DataTransfer) Update(update models.PacketUpdate) {
	dataTransfer.bufMx.Lock()
	defer dataTransfer.bufMx.Unlock()
	dataTransfer.packetBuf[update.ID] = update
}
