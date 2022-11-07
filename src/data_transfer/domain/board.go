package domain

import excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"

type Board struct {
	Name                 string
	PacketTimestampPairs map[uint16]*PacketTimestampPair
}

func NewBoard(rawBoard excelAdapter.BoardDTO) Board {
	return Board{
		Name:                 rawBoard.Name,
		PacketTimestampPairs: NewPacketTimestampPairs(rawBoard.GetPackets()),
	}
}
