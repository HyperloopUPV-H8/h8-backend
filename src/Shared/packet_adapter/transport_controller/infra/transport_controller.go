package infra

import (
	"net"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/transport_controller/infra/sniffer"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/transport_controller/infra/tcp"
	"github.com/google/gopacket"
)

type Port = uint16
type IP = string
type Payload = []byte

type TransportController struct {
	sniffer *sniffer.Sniffer
	server  tcp.Server
}

func NewTransportController(config Config) *TransportController {
	if config.TCPConfig == nil {
		*config.TCPConfig = *tcp.DefaultConfig()
	}

	if config.SnifferConfig == nil {
		*config.SnifferConfig = sniffer.DefaultConfig(config.TCPConfig.RemoteIPs, append(config.TCPConfig.RemoteIPs, getLocalIPs()...))
	}

	return &TransportController{
		sniffer: sniffer.New(config.Device, config.Live, config.SnifferConfig),
		server:  tcp.Open(config.TCPConfig),
	}
}

func (controller TransportController) ReceiveData() ([]byte, gopacket.CaptureInfo) {
	return controller.sniffer.GetNext()
}

func (controller TransportController) OnRead(action func([]byte)) {
	controller.server.SetOnRead(action)
}

func (controller TransportController) Send(addr string, payload []byte) {
	controller.server.Send(addr, payload)
}

func (controller TransportController) Close() {
	controller.server.Close()
	controller.sniffer.Close()
}

func getLocalIPs() []IP {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	ips := make([]IP, 0, len(ifaces))
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ips = append(ips, IP(addr.String()))
		}
	}

	return ips
}
