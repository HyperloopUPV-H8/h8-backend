package message

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/packet"
)

type Payload struct {
	Data fmt.Stringer
	raw  []byte
}

func (message Payload) Kind() packet.Kind {
	return packet.Message
}

func (message Payload) Raw() []byte {
	return message.raw
}

type BlcuAck struct {
	raw []byte
}

func (message BlcuAck) String() string {
	return string(message.raw)
}
