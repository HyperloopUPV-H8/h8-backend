package packetParser

import (
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	packetParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra/dto"
	transportController "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/transport_controller"
)

type PacketAdapter struct {
	controller transportController.TransportController
	parser     packetParser.PacketParser
}

func New(ips []string, packets []excelAdapter.PacketDTO) PacketAdapter {
	return PacketAdapter{
		controller: transportController.NewTransportController(ips),
		parser:     packetParser.NewParser(packets),
	}
}

func (adapter PacketAdapter) ReceiveData() dto.PacketUpdate {
	payload := adapter.controller.ReceiveData()
	decodedPayload := adapter.parser.Decode(payload)
	return decodedPayload
}
