package sniffer

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const SNAPLEN = 1500

type Sniffer struct {
	source *pcap.Handle
	filter string
}

func New(dev string, filter string) (*Sniffer, error) {
	source, err := obtainSource(dev, filter)
	if err != nil {
		return nil, err
	}
	return &Sniffer{
		source: source,
		filter: filter,
	}, nil
}

func obtainSource(dev string, filter string) (*pcap.Handle, error) {
	source, err := pcap.OpenLive(dev, SNAPLEN, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	if err := source.SetBPFFilter(filter); err != nil {
		return nil, err
	}

	return source, nil
}

func (sniffer *Sniffer) Listen(output chan<- []byte) <-chan error {
	errorChan := make(chan error)

	go sniffer.read(output, errorChan)

	return errorChan
}

func (sniffer *Sniffer) read(output chan<- []byte, errorChan chan<- error) {
	for {
		raw, _, err := sniffer.source.ReadPacketData()
		if err != nil {
			errorChan <- err
			close(errorChan)
			return
		}

		output <- gopacket.NewPacket(raw, layers.LayerTypeEthernet, gopacket.DecodeOptions{
			Lazy:   true,
			NoCopy: true,
		}).ApplicationLayer().Payload()
	}
}
