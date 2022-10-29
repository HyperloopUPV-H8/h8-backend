package sniffer

import (
	"strings"

	"github.com/google/gopacket"
)

type IP = string

// The *name*er syntax is standard in go for interfaces that only provide the method *name*,
// even if it's not valid english
type Filterer interface {
	Filter(gopacket.Packet) bool
}

type SourceIPFilter struct {
	SrcIPs []IP
}

func (filter SourceIPFilter) Filter(packet gopacket.Packet) bool {
	srcIP := getPacketSrcIP(packet)

	for _, ip := range filter.SrcIPs {
		if ip == srcIP {
			return true
		}
	}

	return false
}

type DestinationIPFilter struct {
	DstIPs []IP
}

func (filter DestinationIPFilter) Filter(packet gopacket.Packet) bool {
	dstIP := getPacketDstIP(packet)
	for _, ip := range filter.DstIPs {
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
