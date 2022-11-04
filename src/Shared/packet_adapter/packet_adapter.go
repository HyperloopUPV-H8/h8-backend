package packetParser

import (
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	packetParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser"
	packetParserDomain "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/domain"
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

func (adapter PacketAdapter) ReceiveData() packetParserDomain.PacketUpdate {
	payload := adapter.controller.ReceiveData()
	return adapter.parser.Decode(payload)
}
