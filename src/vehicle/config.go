package vehicle

import (
	"github.com/HyperloopUPV-H8/Backend-H8/sniffer"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/parsers"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/parsers/packet_parser"
)

type Config struct {
	Network      NetworkConfig                  `toml:"network"`
	PacketParser packet_parser.Config           `toml:"packet_parser"`
	Protections  parsers.ProtectionParserConfig `toml:"protections"`
}

type NetworkConfig struct {
	TcpClientTag string         `toml:"tcp_client_tag"`
	TcpServerTag string         `toml:"tcp_server_tag"`
	UdpTag       string         `toml:"udp_tag"`
	Sniffer      sniffer.Config `toml:"sniffer"`
}
