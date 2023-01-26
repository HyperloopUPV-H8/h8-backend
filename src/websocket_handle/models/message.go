package models

import "encoding/json"

type Message struct {
	Kind string          `json:"type"`
	Msg  json.RawMessage `json:"msg"`
}

type MessageTarget struct {
	Target []string
	Msg    Message
}
