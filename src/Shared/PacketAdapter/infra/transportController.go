package PacketAdapter

type TransportController struct {
	sniffer Sniffer
	server  Server
}

var (
	snifferTarget string = "\\Device\\NPF_Loopback"
	snifferLive   bool   = true

	serverPort Port = 6000
)

func NewTransportController(validAddrs []string) *TransportController {
	validAddrIPs := stringsToIPs(validAddrs)

	return &TransportController{
		sniffer: NewSniffer(snifferTarget, snifferLive, createFilters(validAddrIPs)),
		server:  OpenServer(serverPort, validAddrIPs),
	}
}

func (controller *TransportController) ReceiveData() []byte {
	return controller.sniffer.GetNextPacket()
}

func (controller *TransportController) ReceiveMessage() []byte {
	return controller.server.Receive()
}

func (controller *TransportController) Send(addr string, payload []byte) {
	controller.server.Send(IP(addr), payload)
}

func (controller *TransportController) AliveConnections() []string {
	return ipsToStrings(controller.server.ConnectedAddresses())
}

func createFilters(validAddrIPs []IP) []Filterer {
	return []Filterer{NewDesiredEndpointsFilter(validAddrIPs)}
}
