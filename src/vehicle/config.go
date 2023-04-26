package vehicle

import (
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/packet_parser"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/protection_parser"
)

type Config struct {
	Network      NetworkConfig            `toml:"network"`
	PacketParser packet_parser.Config     `toml:"packet_parser"`
	Protections  protection_parser.Config `toml:"protections"`
}

type NetworkConfig struct {
	TcpClientTag string `toml:"tcp_client_tag"`
	TcpServerTag string `toml:"tcp_server_tag"`
	UdpTag       string `toml:"udp_tag"`
	Mtu          uint   `toml:"mtu"`
	Interface    string `toml:"interface"`
}
