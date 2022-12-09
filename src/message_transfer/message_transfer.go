package message_transfer

import (
	dataTransferModels "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/models"
	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer/models"
	"github.com/gorilla/websocket"
)

type MessageTransfer struct {
	routines []chan<- dataTransferModels.PacketUpdate
}

func (messageTransfer *MessageTransfer) HandleConn(socket *websocket.Conn) {
	updates := make(chan dataTransferModels.PacketUpdate)
	messageTransfer.routines = append(messageTransfer.routines, updates)

	go func(socket *websocket.Conn, updates <-chan dataTransferModels.PacketUpdate) {
		defer socket.Close()
		for {
			update := <-updates
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

			if err := socket.WriteJSON(message); err != nil {
				return
			}
		}
	}(socket, updates)
}

func (messageTransfer *MessageTransfer) Broadcast(update dataTransferModels.PacketUpdate) {
	for _, routine := range messageTransfer.routines {
		routine <- update
	}
}
