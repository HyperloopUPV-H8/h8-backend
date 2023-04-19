package data

import "github.com/HyperloopUPV-H8/Backend-H8/packet"

type Payload struct {
	Values map[string]packet.Value
	raw    []byte
}

func (data Payload) Kind() packet.Kind {
	return packet.Data
}

func (data Payload) Raw() []byte {
	return data.raw
}
