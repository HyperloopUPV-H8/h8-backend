package data

import "github.com/HyperloopUPV-H8/Backend-H8/packet"

type Payload struct {
	Values map[string]packet.Value
}

func (data Payload) Kind() packet.Kind {
	return packet.Data
}
