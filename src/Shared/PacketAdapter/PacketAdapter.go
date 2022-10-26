package packetadapter

type PacketAdapter struct {
	packetParser        packetparser.PacketParser
	transportController packetparser.TransportController
}

func (pa *PacketAdapter) SendOrder(order orders.OrderDTO) {
	encodedOrder := pa.packetParser.GetEncodedOrder(order)
	pa.transportController.sendTCP(encodedOrder.ip, encodedOrder.bytes)
}

func (pa *PacketAdapter) GetPacketUpdates() []packetparser.PacketUpdate {
	bytesArr := pa.transportController.getPackets()
	packetUpdates := make([]packetparser.PacketUpdate, len(bytesArr))
	for index, bytes := range bytesArr {
		packetUpdates[index] = pa.packetParser.toPacketUpdate(bytes)
	}

	return packetUpdates
}
