package message_transfer

import (
	"log"

	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer/models"
	board_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	ws_models "github.com/HyperloopUPV-H8/Backend-H8/websocket_handle/models"
)

type MessageTransfer struct {
	channel chan ws_models.MessageTarget
}

func New() (*MessageTransfer, chan ws_models.MessageTarget) {
	channel := make(chan ws_models.MessageTarget)
	return &MessageTransfer{channel}, channel
}

func (messageTransfer *MessageTransfer) Broadcast(update board_models.Update) {
	message, err := ws_models.NewMessageTarget([]string{}, "message/update", getMessage(update))
	if err != nil {
		log.Printf("messageTransfer: Broadcast: %s\n", err)
		return
	}
	messageTransfer.channel <- message
}

func getMessage(update board_models.Update) models.Message {
	var message models.Message
	if msg, ok := update.Fields["warning"]; ok {
		message = models.Message{
			ID:          update.ID,
			Description: msg.(string),
			Type:        "warning",
		}
	} else if msg, ok = update.Fields["fault"]; ok {
		message = models.Message{
			ID:          update.ID,
			Description: msg.(string),
			Type:        "fault",
		}
	} else {
		log.Fatalln("MessageTransfer: getMessage: get msg: update does not contain field \"fault\" or \"warning\"")
	}
	return message
}
