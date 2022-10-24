package PacketAdapter

import (
	"net"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/infra/aliases"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/infra/sniffer"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/infra/tcp"
)

type TransportController struct {
	sniffer sniffer.Sniffer
	server  tcp.Server
}

var (
	snifferTarget string = "\\Device\\NPF_Loopback"
	snifferLive   bool   = true

	serverPort aliases.Port = 6000
)

func NewTransportController(validAddrs []string) TransportController {
	validAddrIPs := aliases.StringsToIPs(validAddrs)

	return TransportController{
		sniffer: sniffer.New(snifferTarget, snifferLive, createFilters(validAddrIPs)),
		server:  tcp.Open(serverPort, validAddrIPs),
	}
}

func (controller TransportController) ReceiveData() []byte {
	return controller.sniffer.GetNextValidPayload()
}

func (controller TransportController) ReceiveMessage() [][]byte {
	return aliases.PayloadsToBytes(controller.server.ReceiveNext())
}

func (controller TransportController) Send(addr string, payload []byte) {
	controller.server.Send(aliases.IP(addr), payload)
}

func (controller TransportController) AliveConnections() []string {
	return aliases.IPsToStrings(controller.server.ConnectedAddresses())
}

func (controller TransportController) Close() {
	controller.server.Close()
}

func createFilters(validAddrIPs []aliases.IP) []sniffer.Filterer {
	ipRange := append(validAddrIPs, getLocalIPs()...)
	return []sniffer.Filterer{sniffer.UDPFilter{}, sniffer.SourceIPFilter{SrcIPs: validAddrIPs}, sniffer.DestinationIPFilter{DstIPs: ipRange}}
}

func getLocalIPs() []aliases.IP {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	ips := make([]aliases.IP, 0, len(ifaces))
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ips = append(ips, aliases.IP(addr.String()))
		}
	}

	return ips
}
