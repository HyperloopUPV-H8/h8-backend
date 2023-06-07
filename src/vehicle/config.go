package vehicle

import (
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/packet_parser"
)

type Config struct {
	Boards       []string             `toml:"boards,omitempty"`
	Network      NetworkConfig        `toml:"network"`
	PacketParser packet_parser.Config `toml:"packet_parser"`
	Messages     MessageConfig        `toml:"messages"`
}

type NetworkConfig struct {
	TcpClientTag string `toml:"tcp_client_tag"`
	TcpServerTag string `toml:"tcp_server_tag"`
	UdpTag       string `toml:"udp_tag"`
	Mtu          uint   `toml:"mtu"`
	Interface    string `toml:"interface"`
}

type MessageConfig struct {
	InfoIdKey        string `toml:"info_id_key"`
	WarningIdKey     string `toml:"warning_id_key"`
	FaultIdKey       string `toml:"fault_id_key"`
	ErrorIdKey       string `toml:"error_id_key"`
	BlcuAckId        string `toml:"blcu_ack_id_key"`
	StateOrdersIdKey string `toml:"state_orders_id_key"`
}
