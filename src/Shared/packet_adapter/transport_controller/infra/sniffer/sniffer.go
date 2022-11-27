package sniffer

import (
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type Sniffer struct {
	source *pcap.Handle
}

func New(target string, live bool, config *Config) *Sniffer {
	return &Sniffer{
		source: obtainSource(target, live, config.snaplen, config.promisc, config.timeout, config.bpf),
	}
}

func obtainSource(target string, live bool, snaplen int32, promisc bool, timeout time.Duration, bpf string) *pcap.Handle {
	var (
		handle *pcap.Handle
		err    error
	)

	if live {
		handle, err = pcap.OpenLive(target, snaplen, promisc, timeout)
	} else {
		handle, err = pcap.OpenOffline(target)
	}

	if err != nil {
		log.Fatalf("sniffer: obtain source: %s\n", err)
	}

	if err := handle.SetBPFFilter(bpf); err != nil {
		log.Fatalf("sniffer: obtain source: %s\n", err)
	}

	return handle
}

func (sniffer *Sniffer) GetNext() ([]byte, gopacket.CaptureInfo) {
	payload, info, err := sniffer.source.ReadPacketData()
	if err != nil {
		log.Fatalf("sniffer: get next: %s\n", err)
	}

	return payload, info
}

func (sniffer *Sniffer) Close() {
	sniffer.source.Close()
}
