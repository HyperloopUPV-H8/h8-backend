package domain

import excelParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain/board"

type Board struct {
	Name                 string
	PacketTimestampPairs map[uint16]PacketTimestampPair
}

func NewBoard(rawBoard excelParser.Board) Board {
	return Board{
		Name:                 rawBoard.Name,
		PacketTimestampPairs: NewPacketTimestampPairs(rawBoard.GetPackets()),
	}
}
