package domain

import (
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	packetParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/domain"
)

type PodData struct {
	Boards map[string]Board
}

func (podData *PodData) UpdatePacket(pu packetParser.PacketUpdate) {
	packet := podData.GetPacket(pu.ID)
	packet.UpdatePacket(pu)
}

func (podData *PodData) GetPacket(id uint16) *PacketTimestampPair {
	for _, board := range podData.Boards {
		packetTimeStampPair, exists := board.PacketTimestampPairs[id]
		if exists {
			return &packetTimeStampPair
		}
	}
	return nil
}

func NewPodData(rawBoards map[string]excelAdapter.BoardDTO) PodData {
	return PodData{
		Boards: getBoards(rawBoards),
	}
}

func getBoards(rawBoards map[string]excelAdapter.BoardDTO) map[string]Board {
	boards := make(map[string]Board)
	for name, board := range rawBoards {
		board := Board{
			Name:                 name,
			PacketTimestampPairs: NewPacketTimestampPairs(board.GetPackets()),
		}
		boards[board.Name] = board
	}
	return boards
}
