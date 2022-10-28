package sniffer

import (
	"strings"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/infra/aliases"
	"github.com/google/gopacket"
)

// The *name*er syntax is standard in go for interfaces that only provide the method *name*,
// even if it's not valid english
type Filterer interface {
	Filter(gopacket.Packet) bool
}

type SourceIPFilter struct {
	SrcIPs []aliases.IP
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
	DstIPs []aliases.IP
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

func getPacketSrcIP(packet gopacket.Packet) (src aliases.IP) {
	networkLayer := packet.NetworkLayer()
	if networkLayer == nil {
		return
	}

	netFlow := networkLayer.NetworkFlow()

	return aliases.IP(netFlow.Src().String())
}

func getPacketDstIP(packet gopacket.Packet) (src aliases.IP) {
	networkLayer := packet.NetworkLayer()
	if networkLayer == nil {
		return
	}

	netFlow := networkLayer.NetworkFlow()

	return aliases.IP(netFlow.Src().String())
}
