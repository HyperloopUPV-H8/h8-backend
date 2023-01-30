package message_transfer

import (
	"log"

	dataTransferModels "github.com/HyperloopUPV-H8/Backend-H8/data_transfer/models"
	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer/models"
	ws_models "github.com/HyperloopUPV-H8/Backend-H8/websocket_handle/models"
)

type MessageTransfer struct {
	channel chan ws_models.MessageTarget
}

func New() (*MessageTransfer, chan ws_models.MessageTarget) {
	channel := make(chan ws_models.MessageTarget)
	return &MessageTransfer{channel}, channel
}

func (messageTransfer *MessageTransfer) Broadcast(update dataTransferModels.PacketUpdate) {
	message, err := ws_models.NewMessageTarget([]string{}, "message/update", getMessage(update))
	if err != nil {
		log.Printf("messageTransfer: Broadcast: %s\n", err)
		return
	}
	messageTransfer.channel <- message
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
			Type:        "fault",
		}
	}
	return message
}
