package infra

type TransportController struct {
	sniffer          *Sniffer
	aliveConnections map[string]*Connection
	deadConnections  map[string]*Ping
	addresses        []string
}

func (controller *TransportController) Receive() []byte {
	// Missing parts
	return controller.sniffer.GetNextPacket()
}

func (controller *TransportController) Send(addr string, payload []byte) {
	// TODO
}

func (controller *TransportController) AliveConnections() []string {
	// TODO
	return nil
}
