package connection_transfer

import (
	"encoding/json"
	"log"

	"github.com/HyperloopUPV-H8/Backend-H8/connection_transfer/models"
	"github.com/gorilla/websocket"
	"github.com/kjk/betterguid"
)

type ConnectionTransfer struct {
	boardStatus map[string]models.Connection
	sockets     map[string]*websocket.Conn
}

func (connectionTransfer *ConnectionTransfer) HandleConn(socket *websocket.Conn) {
	connectionTransfer.sockets[betterguid.New()] = socket
}

func (connectionTransfer *ConnectionTransfer) Update(name string, up bool) {
	connectionTransfer.boardStatus[name] = models.Connection{
		Name:        name,
		IsConnected: up,
	}

	message, err := json.Marshal(mapToArray(connectionTransfer.boardStatus))
	if err != nil {
		log.Fatalf("connection transfer: Update: %s\n", err)
	}

	for id, socket := range connectionTransfer.sockets {
		if err := socket.WriteMessage(websocket.TextMessage, message); err != nil {
			socket.Close()
			delete(connectionTransfer.sockets, id)
		}
	}
}

func (ConnectionTransfer *ConnectionTransfer) Close() {
	for _, socket := range ConnectionTransfer.sockets {
		socket.Close()
	}
}

func mapToArray(input map[string]models.Connection) []models.Connection {
	output := make([]models.Connection, 0, len(input))
	for _, value := range input {
		output = append(output, value)
	}
	return output
}
