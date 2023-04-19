package message

import (
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/packet"
)

type Payload struct {
	Data fmt.Stringer
}

func (message Payload) Kind() packet.Kind {
	return packet.Message
}

type BlcuAck struct {
	raw []byte
}

func (message BlcuAck) String() string {
	return string(message.raw)
}
