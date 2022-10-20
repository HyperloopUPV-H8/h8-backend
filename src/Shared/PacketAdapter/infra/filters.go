package infra

import "github.com/google/gopacket"

// The *name*er syntax is standard in go for interfaces that only provide the method *name*,
// even if it's not valid english
type Filterer interface {
	Filter(gopacket.Packet) bool
}

type DifferentEndpointIPFilter struct {
	endpointIP IP
}

func NewDifferentEndpointIPFilter(endpointIP IP) DifferentEndpointIPFilter {
	return DifferentEndpointIPFilter{endpointIP: endpointIP}
}

func (filter DifferentEndpointIPFilter) Filter(packet gopacket.Packet) bool {
	networkL := packet.NetworkLayer()
	flow := networkL.NetworkFlow()
	return IP(flow.Dst().String()) != filter.endpointIP && IP(flow.Src().String()) != filter.endpointIP
}
