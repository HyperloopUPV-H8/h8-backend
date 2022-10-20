package infra

import (
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

const (
	liveHandleSnaplen int32         = 2147483647 // max int32 value
	liveHandlePromisc bool          = true
	liveHandleTimeout time.Duration = pcap.BlockForever
)

type Sniffer struct {
	source  *gopacket.PacketSource
	filters []Filterer
}

func NewSniffer(target string, live bool, filters []Filterer) Sniffer {
	return Sniffer{
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
		if !filter.Filter(*packet) {
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

		if !sniffer.applyFilters(&packet) {
			continue
		}

		applicationLayer := packet.ApplicationLayer()

		if applicationLayer == nil {
			continue
		}

		return applicationLayer.Payload()
	}
}
