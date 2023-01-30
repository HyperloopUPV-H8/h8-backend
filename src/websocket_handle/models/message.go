package models

import "encoding/json"

type Message struct {
	Topic string          `json:"topic"`
	Msg   json.RawMessage `json:"msg"`
}

type MessageTarget struct {
	Target []string
	Msg    Message
}

func NewMessageTarget(target []string, topic string, msg any) (MessageTarget, error) {
	msgRaw, err := json.Marshal(msg)
	return MessageTarget{
		Target: target,
		Msg: Message{
			Topic: topic,
			Msg:   msgRaw,
		},
	}, err
}

func NewMessageTargetRaw(target []string, topic string, msg json.RawMessage) MessageTarget {
	return MessageTarget{
		Target: target,
		Msg: Message{
			Topic: topic,
			Msg:   msg,
		},
	}
}
