package models

type Message struct {
	Type string `json:"type"`
	Msg  any    `json:"msg"`
}

type MessageTarget struct {
	Target []string
	Msg    Message
}
