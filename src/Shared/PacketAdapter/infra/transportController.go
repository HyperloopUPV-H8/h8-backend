package PacketAdapter

import "net"

type TransportController struct {
	sniffer Sniffer
	server  Server
}

var (
	snifferTarget string = "\\Device\\NPF_Loopback"
	snifferLive   bool   = true

	serverPort Port = 6100
)

func NewTransportController(validAddrs []string) TransportController {
	validAddrIPs := stringsToIPs(validAddrs)

	return TransportController{
		sniffer: NewSniffer(snifferTarget, snifferLive, createFilters(validAddrIPs)),
		server:  OpenServer(serverPort, validAddrIPs),
	}
}

func (controller TransportController) ReceiveData() []byte {
	return controller.sniffer.GetNextValidPayload()
}

func (controller TransportController) ReceiveMessage() [][]byte {
	return payloadsToBytes(controller.server.ReceiveNext())
}

func (controller TransportController) Send(addr string, payload []byte) {
	controller.server.Send(IP(addr), payload)
}

func (controller TransportController) AliveConnections() []string {
	return ipsToStrings(controller.server.ConnectedAddresses())
}

func (controller TransportController) Close() {
	controller.server.Close()
}

func createFilters(validAddrIPs []IP) []Filterer {
	ipRange := append(validAddrIPs, getLocalIPs()...)
	return []Filterer{UDPFilter{}, SourceIPFilter{validAddrIPs}, DestinationIPFilter{ipRange}}
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
