package message_transfer

import (
	"encoding/json"
	"errors"

	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer/models"
)

const (
	MESSAGE_TRANSFER_NAME  = "messageTransfer"
	MESSAGE_TRANSFER_TOPIC = "message/update"
)

var (
	messageTransfer *MessageTransfer
)

func Get() *MessageTransfer {
	if messageTransfer == nil {
		initMessageTransfer()
	}

	return messageTransfer
}

func initMessageTransfer() {
	messageTransfer = &MessageTransfer{defaultSendMessage}
}

type MessageTransfer struct {
	sendMessage func(topic string, payload any, targets ...string) error
}

func (messageTransfer *MessageTransfer) SendMessage(message models.Message) error {
	return messageTransfer.sendMessage(MESSAGE_TRANSFER_TOPIC, message)
}

func (messageTransfer *MessageTransfer) UpdateMessage(topic string, payload json.RawMessage, source string) {
}

func (messageTransfer *MessageTransfer) SetSendMessage(sendMessage func(topic string, payload any, targets ...string) error) {
	messageTransfer.sendMessage = sendMessage
}

func (messageTransfer *MessageTransfer) HandlerName() string {
	return MESSAGE_TRANSFER_NAME
}

func defaultSendMessage(string, any, ...string) error {
	return errors.New("message transfer must be registered before using")
}
