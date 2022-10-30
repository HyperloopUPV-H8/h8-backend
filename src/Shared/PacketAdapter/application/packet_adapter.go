package application

import (
	excelParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain/board"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/infra"
)

type PacketAdapter struct {
	controller infra.TransportController
	parser     domain.PacketParser
}

func New(ips []string, packets []excelParser.Packet) PacketAdapter {
	return PacketAdapter{
		controller: infra.NewTransportController(ips),
		parser:     domain.NewParser(packets),
	}
}

func (adapter PacketAdapter) ReadData() domain.PacketUpdate {
	payload := adapter.controller.ReceiveData()
	return adapter.parser.Decode(payload)
}
