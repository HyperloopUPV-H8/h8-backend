package message_transfer

import (
	dataTransferModels "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/models"
	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer/models"
	"github.com/gorilla/websocket"
	"github.com/kjk/betterguid"
)

type MessageTransfer struct {
	sockets map[string]*websocket.Conn
}

func (messageTransfer *MessageTransfer) HandleConn(socket *websocket.Conn) {
	messageTransfer.sockets[betterguid.New()] = socket
}

func (messageTransfer *MessageTransfer) Broadcast(update dataTransferModels.PacketUpdate) {
	message := getMessage(update)
	for id, socket := range messageTransfer.sockets {
		if err := socket.WriteJSON(message); err != nil {
			socket.Close()
			delete(messageTransfer.sockets, id)
		}
	}
}

func (messageTransfer *MessageTransfer) Close() {
	for _, socket := range messageTransfer.sockets {
		socket.Close()
	}
}

func getMessage(update dataTransferModels.PacketUpdate) models.Message {
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
	return message
}
