package models

import "encoding/json"

type Message struct {
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload"`
}

func NewMessage(topic string, message any) (Message, error) {
	messageRaw, err := json.Marshal(message)

	return Message{
		Topic:   topic,
		Payload: messageRaw,
	}, err
}
