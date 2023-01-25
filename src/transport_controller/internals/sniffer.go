package internals

import (
	"log"

	"github.com/HyperloopUPV-H8/Backend-H8/transport_controller/models"
	"github.com/google/gopacket/pcap"
)

const etherType = 12
const ipType = etherType + 11
const typeIPIP = 0x04
const typeUDP = 0x11
const typeTCP = 0x06
const udpOffset = 42
const tcpOffset = 54

type Sniffer struct {
	source *pcap.Handle
	config models.Config
}

func OpenSniffer(device string, live bool, config models.Config) *Sniffer {
	sniffer := &Sniffer{
		source: obtainSource(device, live, config),
		config: config,
	}

	go sniffer.StartReading()

	return sniffer
}

func obtainSource(device string, live bool, config models.Config) *pcap.Handle {
	var (
		source *pcap.Handle
		err    error
	)

	log.Println(pcap.FindAllDevs())

	if live {
		source, err = pcap.OpenLive(device, config.Snaplen, config.Promisc, config.Timeout)
	} else {
		source, err = pcap.OpenOffline(device)
	}

	if err != nil {
		log.Fatalf("sniffer: obtainSource: %s\n", err)
	}

	err = source.SetBPFFilter(config.BPF)
	if err != nil {
		log.Fatalf("sniffer: obtainSource: %s\n", err)
	}

	return source
}

func (sniffer *Sniffer) StartReading() {
	for {
		payload, _, err := sniffer.source.ReadPacketData()
		if err != nil {
			continue
		}

		if payload[etherType] != 0x08 || payload[etherType+1] != 0x00 {
			continue
		}

		switch payload[ipType] {
		case typeIPIP:
			switch payload[ipType+20] {
			case typeUDP:
				sniffer.config.Dump <- payload[udpOffset+20:]
			case typeTCP:
				sniffer.config.Dump <- payload[tcpOffset+20:]
			}
		case typeUDP:
			sniffer.config.Dump <- payload[udpOffset:]
		case typeTCP:
			sniffer.config.Dump <- payload[tcpOffset:]
		}

	}
}

func (sniffer *Sniffer) Close() {
	sniffer.source.Close()
}

func (sniffer *Sniffer) Stats() (recieved int, dropped int) {
	stats, err := sniffer.source.Stats()
	if err != nil {
		return
	}
	return stats.PacketsReceived, stats.PacketsDropped
}
