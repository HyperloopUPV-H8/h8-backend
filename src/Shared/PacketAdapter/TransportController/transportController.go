package transportController

import (
	"net"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/TransportController/sniffer"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/TransportController/tcp"
)

type Port = uint16
type IP = string
type Payload = []byte

type TransportController struct {
	sniffer sniffer.Sniffer
	server  tcp.Server
}

var (
	snifferTarget string = "\\Device\\NPF_Loopback"
	snifferLive   bool   = true

	serverPort Port = 6000
)

func NewTransportController(validAddrs []string) TransportController {

	return TransportController{
		sniffer: sniffer.New(snifferTarget, snifferLive, createFilters(validAddrs)),
		server:  tcp.Open(serverPort, validAddrs),
	}
}

func (controller TransportController) ReceiveData() []byte {
	return controller.sniffer.GetNextValidPayload()
}

func (controller TransportController) ReceiveMessages() [][]byte {
	return controller.server.ReceiveNext()
}

func (controller TransportController) Send(addr string, payload []byte) {
	controller.server.Send(addr, payload)
}

func (controller TransportController) AliveConnections() []string {
	return controller.server.ConnectedAddresses()
}

func (controller TransportController) Close() {
	controller.server.Close()
}

func createFilters(validAddrIPs []IP) []sniffer.Filterer {
	ipRange := append(validAddrIPs, getLocalIPs()...)
	return []sniffer.Filterer{sniffer.UDPFilter{}, sniffer.SourceIPFilter{SrcIPs: validAddrIPs}, sniffer.DestinationIPFilter{DstIPs: ipRange}}
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
