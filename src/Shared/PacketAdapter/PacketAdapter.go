package packetadapter

import (
	packetparser "github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain/packet_parser"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain/packet_parser/domain"
	transportcontroller "github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/infra"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/excelAdapter/dto"
)

type PacketAdapter struct {
	packetParser        packetparser.PacketParser
	transportController transportcontroller.TransportController
}

func New(packetDTOs []dto.PacketDTO, boardIps []string) PacketAdapter {
	packetAdapter := PacketAdapter{
		packetParser:        packetparser.New(packetDTOs),
		transportController: transportcontroller.NewTransportController(boardIps),
	}

	return packetAdapter
}

func (pa *PacketAdapter) GetPacketUpdate() domain.PacketUpdate {
	bytes := pa.transportController.ReceiveData()
	return pa.packetParser.Decode(bytes)
}

// func (pa *PacketAdapter) SendOrder(order orders.OrderDTO) {
// 	encodedOrder := pa.packetParser.GetEncodedOrder(order)
// 	pa.transportController.sendTCP(encodedOrder.ip, encodedOrder.bytes)
// }
