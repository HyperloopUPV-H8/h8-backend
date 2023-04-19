package vehicle

import (
	"github.com/HyperloopUPV-H8/Backend-H8/packet/parsers"
	"github.com/HyperloopUPV-H8/Backend-H8/sniffer"
)

type Config struct {
	TcpClientTag string `toml:"tcp_client_tag"`
	TcpServerTag string `toml:"tcp_server_tag"`
	UdpTag       string `toml:"udp_tag"`
	Network      sniffer.Config
	Parser       parsers.Config
}
