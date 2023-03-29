package models

import "encoding/json"

type MessageHandler interface {
	UpdateMessage(topic string, payload json.RawMessage, source string)
	SetSendMessage(func(topic string, payload any, targets ...string) error)
	HandlerName() string
}
