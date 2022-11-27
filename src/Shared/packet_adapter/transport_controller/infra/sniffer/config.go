package sniffer

import (
	"time"

	"github.com/google/gopacket/pcap"
)

type Config struct {
	snaplen int32
	promisc bool
	bpf     string
	timeout time.Duration
}

func DefaultConfig(srcAddrs []string, dstAddrs []string) *Config {
	return &Config{
		snaplen: ^int32(0),
		promisc: true,
		timeout: pcap.BlockForever,
		bpf:     getFilters(srcAddrs, dstAddrs),
	}
}
