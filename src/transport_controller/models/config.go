package models

import (
	"net"
	"time"

	"github.com/google/gopacket/pcap"
)

type Config struct {
	Dump         chan []byte
	Snaplen      int32
	Promisc      bool
	Timeout      time.Duration
	BPF          string
	OnConnUpdate func(*net.TCPAddr, bool)
}

var defaultConfig = Config{
	Snaplen:      ^int32(0),
	Promisc:      true,
	Timeout:      pcap.BlockForever,
	Dump:         make(chan []byte, 1024),
	BPF:          "udp",
	OnConnUpdate: func(t *net.TCPAddr, b bool) {},
}
