package models

type Message struct {
	Kind string `json:"kind"`
	Msg  any    `json:"msg"`
}

type MessageTarget struct {
	Target []string
	Msg    Message
}
