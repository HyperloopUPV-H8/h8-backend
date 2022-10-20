package infra

type TransportController struct {
	sniffer Sniffer
	server  Server
}

const (
	snifferTarget string = ""
	snifferLive   bool   = true

	serverPort Port = 6000
)

func New(validAddrs []string) *TransportController {
	return &TransportController{
		sniffer: NewSniffer(snifferTarget, snifferLive, nil),
		server:  OpenServer(serverPort, stringsToIPs(validAddrs)),
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
	return controller.server.ConnectedAddresses()
}
