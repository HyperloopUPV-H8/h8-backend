package packetParser

import (
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelAdapter/domain"
	packetParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/PacketParser"
	packetParserDomain "github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/PacketParser/domain"
	transportController "github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/TransportController"
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
