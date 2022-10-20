package infra

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

const (
	liveHandleSnaplen = 2_147_483_647
	liveHandlePromisc = true
	liveHandleTimeout = pcap.BlockForever
)

type Sniffer struct {
	source  *gopacket.PacketSource
	filters []func(*gopacket.Packet) bool
}

func NewSniffer(target string, live bool, filters []func(*gopacket.Packet) bool) *Sniffer {
	return &Sniffer{
		source:  obtainSource(target, live),
		filters: filters,
	}
}

func obtainSource(target string, live bool) (source *gopacket.PacketSource) {
	var (
		handle *pcap.Handle
		err    error
	)

	if live {
		handle, err = pcap.OpenLive(target, liveHandleSnaplen, liveHandlePromisc, liveHandleTimeout)
	} else {
		handle, err = pcap.OpenOffline(target)
	}

	if err != nil {
		panic(err)
	}

	return gopacket.NewPacketSource(handle, handle.LinkType())
}

func (sniffer *Sniffer) applyFilters(packet *gopacket.Packet) bool {
	for _, filter := range sniffer.filters {
		if !filter(packet) {
			return false
		}
	}

	return true
}

func (sniffer *Sniffer) GetNextPacket() []byte {
	for {
		packet, err := sniffer.source.NextPacket()
		if err != nil {
			continue
		}

		if sniffer.applyFilters(&packet) {
			return packet.ApplicationLayer().Payload()
		}
	}
}
