package models

import (
	"fmt"
	"time"
)

type Message struct {
	Timestamp uint64
	Msg       string
}

func NewMessage(msg string) Message {
	return Message{
		Timestamp: uint64(time.Now().UnixNano()),
		Msg:       msg,
	}
}

func (message *Message) ToCSV() []string {
	return []string{fmt.Sprint(message.Timestamp), message.Msg}
}
