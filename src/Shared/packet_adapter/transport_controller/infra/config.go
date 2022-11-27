package infra

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/transport_controller/infra/sniffer"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/transport_controller/infra/tcp"
)

type Config struct {
	Device        string
	Live          bool
	SnifferConfig *sniffer.Config
	TCPConfig     *tcp.Config
}
