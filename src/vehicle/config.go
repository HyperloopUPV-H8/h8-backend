package vehicle

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/packet_parser"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Boards       []string             `toml:"boards,omitempty"`
	Network      NetworkConfig        `toml:"network"`
	PacketParser packet_parser.Config `toml:"packet_parser"`
	Messages     MessageConfig        `toml:"messages"`
}

type NetworkConfig struct {
	TcpClientTag      string  `toml:"tcp_client_tag"`
	TcpServerTag      string  `toml:"tcp_server_tag"`
	UdpTag            string  `toml:"udp_tag"`
	Mtu               uint    `toml:"mtu"`
	Interface         string  `toml:"interface"`
	KeepaliveInterval *string `toml:"keepalive,omitempty"`
	WriteTimeout      *string `toml:"timeout,omitempty"`
}

func (networkConfig NetworkConfig) GetKeepaliveInterval() *time.Duration {
	if networkConfig.KeepaliveInterval == nil {
		return nil
	}
	interval, err := time.ParseDuration(*networkConfig.KeepaliveInterval)
	if err != nil {
		log.Fatal().Stack().Err(err).Str("interval", *networkConfig.KeepaliveInterval).Msg("error parsing keepalive interval")
	}
	return &interval
}

func (networkConfig NetworkConfig) GetWriteTimeout() *time.Duration {
	if networkConfig.WriteTimeout == nil {
		return nil
	}
	timeout, err := time.ParseDuration(*networkConfig.WriteTimeout)
	if err != nil {
		log.Fatal().Stack().Err(err).Str("timeout", *networkConfig.WriteTimeout).Msg("error parsing write timeout")
	}
	return &timeout
}

type MessageConfig struct {
	InfoIdKey              string `toml:"info_id_key"`
	WarningIdKey           string `toml:"warning_id_key"`
	FaultIdKey             string `toml:"fault_id_key"`
	BlcuAckId              string `toml:"blcu_ack_id_key"`
	AddStateOrdersIdKey    string `toml:"add_state_orders_id_key"`
	RemoveStateOrdersIdKey string `toml:"remove_state_orders_id_key"`
}
