package PacketAdapter

import (
	"strings"

	"github.com/google/gopacket"
)

// The *name*er syntax is standard in go for interfaces that only provide the method *name*,
// even if it's not valid english
type Filterer interface {
	Filter(gopacket.Packet) bool
}

type SourceIPFilter struct {
	srcIPs []IP
}

func (filter SourceIPFilter) Filter(packet gopacket.Packet) bool {
	srcIP := getPacketSrcIP(packet)

	for _, ip := range filter.srcIPs {
		if ip == srcIP {
			return true
		}
	}

	return false
}

type DestinationIPFilter struct {
	dstIPs []IP
}

func (filter DestinationIPFilter) Filter(packet gopacket.Packet) bool {
	dstIP := getPacketDstIP(packet)

	for _, ip := range filter.dstIPs {
		if ip == dstIP {
			return true
		}
	}

	return false
}

type UDPFilter struct{}

func (filter UDPFilter) Filter(packet gopacket.Packet) bool {
	transportLayer := packet.TransportLayer()
	if transportLayer == nil {
		return false
	}

	return strings.HasPrefix(gopacket.LayerString(transportLayer), "UDP")
}

func getPacketSrcIP(packet gopacket.Packet) (src IP) {
	networkLayer := packet.NetworkLayer()
	if networkLayer == nil {
		return
	}

	netFlow := networkLayer.NetworkFlow()

	return IP(netFlow.Src().String())
}

func getPacketDstIP(packet gopacket.Packet) (src IP) {
	networkLayer := packet.NetworkLayer()
	if networkLayer == nil {
		return
	}

	netFlow := networkLayer.NetworkFlow()

	return IP(netFlow.Src().String())
}
