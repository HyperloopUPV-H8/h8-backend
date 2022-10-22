package PacketAdapter

import "github.com/google/gopacket"

// The *name*er syntax is standard in go for interfaces that only provide the method *name*,
// even if it's not valid english
type Filterer interface {
	Filter(gopacket.Packet) bool
}

type DesiredEndpointsFilter struct {
	endpointIPs []IP
}

func (filter DesiredEndpointsFilter) Filter(packet gopacket.Packet) bool {
	srcIP, dstIP := getPacketEndpointIPs(packet)

	srcCheck := false
	dstCheck := false
	for _, ip := range filter.endpointIPs {
		if srcIP == ip {
			srcCheck = true
		} else if dstIP == ip {
			dstCheck = true
		}
	}

	if srcCheck && dstCheck {
		return true
	}

	return false
}

func getPacketEndpointIPs(packet gopacket.Packet) (src, dst IP) {
	networkLayer := packet.NetworkLayer()
	if networkLayer == nil {
		return
	}

	netFlow := networkLayer.NetworkFlow()
	src = IP(netFlow.Src().String())
	dst = IP(netFlow.Dst().String())

	return src, dst
}
