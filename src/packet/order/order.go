package order

import "github.com/HyperloopUPV-H8/Backend-H8/packet"

type Payload struct {
	Values  map[string]packet.Value
	Enabled map[string]bool
	raw     []byte
}

func (order Payload) Kind() packet.Kind {
	return packet.Order
}

func (order Payload) Raw() []byte {
	return order.raw
}
