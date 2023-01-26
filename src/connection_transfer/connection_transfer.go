package connection_transfer

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/HyperloopUPV-H8/Backend-H8/connection_transfer/models"
	ws_models "github.com/HyperloopUPV-H8/Backend-H8/websocket_handle/models"
)

type ConnectionTransfer struct {
	writeMx     sync.Mutex
	boardStatus map[string]models.Connection
	channel     chan ws_models.MessageTarget
}

func New() (*ConnectionTransfer, chan ws_models.MessageTarget) {
	connectionTransfer := &ConnectionTransfer{
		writeMx:     sync.Mutex{},
		boardStatus: make(map[string]models.Connection),
		channel:     make(chan ws_models.MessageTarget),
	}

	go connectionTransfer.run()

	return connectionTransfer, connectionTransfer.channel
}

func (connectionTransfer *ConnectionTransfer) run() {
	for range connectionTransfer.channel {
		connectionTransfer.send()
	}
}

func (connectionTransfer *ConnectionTransfer) Update(name string, up bool) {
	connectionTransfer.writeMx.Lock()
	defer connectionTransfer.writeMx.Unlock()
	connectionTransfer.boardStatus[name] = models.Connection{
		Name:        name,
		IsConnected: up,
	}

	connectionTransfer.send()
}

func (connectionTransfer *ConnectionTransfer) send() {
	msg, err := json.Marshal(mapToArray(connectionTransfer.boardStatus))
	if err != nil {
		log.Printf("connectionTransfer: send: %s\n", err)
		return
	}
	connectionTransfer.channel <- ws_models.MessageTarget{
		Target: []string{},
		Msg: ws_models.Message{
			Kind: "connection/update",
			Msg:  msg,
		},
	}
}

func mapToArray(input map[string]models.Connection) []models.Connection {
	output := make([]models.Connection, 0, len(input))
	for _, value := range input {
		output = append(output, value)
	}
	return output
}
