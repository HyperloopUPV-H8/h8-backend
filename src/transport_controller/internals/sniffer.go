package internals

import (
	"log"
	"os"

	"github.com/HyperloopUPV-H8/Backend-H8/transport_controller/models"
	"github.com/google/gopacket/pcap"
)

const ETHER_TYPE = 12
const IP_TYPE = ETHER_TYPE + 11
const TYPE_IPIP = 0x04
const TYPE_UDP = 0x11
const TYPE_TCP = 0x06
const UDP_OFFSET = 42
const TCP_OFFSET = 54

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

		if os.Getenv("SNIFFER_DEV") == "\\Device\\NPF_Loopback" {
			sniffer.config.Dump <- payload[32:]
		} else {
			if payload[ETHER_TYPE] != 0x08 || payload[ETHER_TYPE+1] != 0x00 {
				continue
			}

			switch payload[IP_TYPE] {
			case TYPE_IPIP:
				switch payload[IP_TYPE+20] {
				case TYPE_UDP:
					sniffer.config.Dump <- payload[UDP_OFFSET+20:]
				case TYPE_TCP:
					sniffer.config.Dump <- payload[TCP_OFFSET+20:]
				}
			case TYPE_UDP:
				sniffer.config.Dump <- payload[UDP_OFFSET:]
			case TYPE_TCP:
				sniffer.config.Dump <- payload[TCP_OFFSET:]
			}
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
