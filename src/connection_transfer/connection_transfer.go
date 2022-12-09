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
			status := make([]models.Connection, 0, len(connHandle.BoardStatus))
			for _, sts := range connHandle.BoardStatus {
				status = append(status, sts)
			}
			socket.WriteJSON(status)
		}
	}(socket, connHandle.Updates)
}
