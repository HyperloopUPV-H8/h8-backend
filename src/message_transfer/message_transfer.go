package message_transfer

import (
	dataTransferModels "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/models"
	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer/models"
	"github.com/gorilla/websocket"
)

type MessageTransfer struct {
	sockets map[string]*websocket.Conn
}

func (messageTransfer *MessageTransfer) HandleConn(socket *websocket.Conn) {
	messageTransfer.sockets[socket.RemoteAddr().String()] = socket
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

	closed := make([]string, 0, len(messageTransfer.sockets))
	for name, socket := range messageTransfer.sockets {
		if err := socket.WriteJSON(message); err != nil {
			closed = append(closed, name)
		}
	}
	for _, name := range closed {
		delete(messageTransfer.sockets, name)
	}
}
