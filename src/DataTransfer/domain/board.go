package domain

import "github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/application/interfaces"

type Board struct {
	Name    string
	Packets map[uint16]Packet
}

func NewBoard(rawBoard interfaces.Board) Board {
	return Board{
		Name:    rawBoard.Name(),
		Packets: NewPackets(rawBoard.GetPackets()),
	}
}
