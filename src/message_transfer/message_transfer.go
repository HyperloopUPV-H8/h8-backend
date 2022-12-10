package message_transfer

import (
	"time"

	dataTransferModels "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/models"
	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer/models"
	"github.com/gorilla/websocket"
)

type MessageTransfer struct {
	sockets map[int64]*websocket.Conn
}

func (messageTransfer *MessageTransfer) HandleConn(socket *websocket.Conn) {
	messageTransfer.sockets[time.Now().UnixNano()] = socket
}

func (messageTransfer *MessageTransfer) Broadcast(update dataTransferModels.PacketUpdate) {
	var message models.Message
	if msg, ok := update.Values["warning"]; ok {
		message = models.Message{
			ID:          update.ID,
			Description: msg,
			Type:        "warning",
		}
	} else if msg, ok = update.Values["fault"]; ok {
		message = models.Message{
			ID:          update.ID,
			Description: msg,
			Type:        "warning",
		}
	}

	closed := make([]int64, 0, len(messageTransfer.sockets))
	for id, socket := range messageTransfer.sockets {
		if err := socket.WriteJSON(message); err != nil {
			closed = append(closed, id)
		}
	}
	for _, id := range closed {
		delete(messageTransfer.sockets, id)
	}
}
