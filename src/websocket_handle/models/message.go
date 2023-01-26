package models

type Message struct {
	Kind string `json:"type"`
	Msg  any    `json:"msg"`
}

type MessageTarget struct {
	Target []string
	Msg    Message
}
