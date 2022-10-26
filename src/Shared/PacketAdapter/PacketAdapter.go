package packetadapter

type PacketAdapter struct {
	packetParser        packetparser.PacketParser
	transportController packetparser.TransportController
}

func (pa *PacketAdapter) sendOrder(order orders.OrderDTO) {
	buf, ip := pa.packetParser.GetOrderBufAndAddress(order)
	pa.transportController.sendTCP(ip, bytes)
}

func (pa *PacketAdapter) getUpdates() []packetparser.PacketUpdate {
	bytesArr := pa.transportController.getPackets()
	packetUpdates := make([]packetparser.PacketUpdate, len(bytesArr))
	for index, bytes := range bytesArr {
		packetUpdates[index] = pa.packetParser.toPacketUpdate(bytes)
	}

	return packetUpdates
}
