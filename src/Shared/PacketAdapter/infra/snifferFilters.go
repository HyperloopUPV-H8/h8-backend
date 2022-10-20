package infra

import "github.com/google/gopacket"

// The *name*er syntax is standard in go for interfaces that only provide the method *name*,
// even if it's not valid english
type Filterer interface {
	Filter(gopacket.Packet) bool
}

type DesiredEndpointsFilter struct {
	endpointIPs []IP
}

func NewDesiredEndpointsFilter(endpointIPs []IP) DesiredEndpointsFilter {
	return DesiredEndpointsFilter{endpointIPs: endpointIPs}
}

func (filter DesiredEndpointsFilter) Filter(packet gopacket.Packet) bool {
	networkLayer := packet.NetworkLayer()
	if networkLayer == nil {
		return false
	}

	netFlow := networkLayer.NetworkFlow()
	srcIP := IP(netFlow.Src().String())
	dstIP := IP(netFlow.Dst().String())

	srcCheck := false
	dstCheck := false
	for _, ip := range filter.endpointIPs {
		if srcIP == ip {
			srcCheck = true
		} else if dstIP == ip {
			dstCheck = true
		}

		if srcCheck && dstCheck {
			return true
		}
	}

	return false
}
