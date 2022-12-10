package connection_transfer

import (
	"github.com/HyperloopUPV-H8/Backend-H8/connection_transfer/models"
	"github.com/gorilla/websocket"
)

type ConnectionTransfer struct {
	BoardStatus map[string]models.Connection
	Updates     <-chan models.Connection
}

func (connHandle *ConnectionTransfer) HandleConn(socket *websocket.Conn) {
	go func(socket *websocket.Conn, updates <-chan models.Connection) {
		for update := range updates {
			connHandle.BoardStatus[update.Name] = update
			socket.WriteJSON(mapToArray(connHandle.BoardStatus))
		}
	}(socket, connHandle.Updates)
}

func mapToArray(input map[string]models.Connection) []models.Connection {
	output := make([]models.Connection, 0, len(input))
	for _, value := range input {
		output = append(output, value)
	}
	return output
}
