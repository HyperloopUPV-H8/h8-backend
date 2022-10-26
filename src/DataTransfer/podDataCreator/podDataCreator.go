package podDataCreator

import (
	domain "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/excelAdapter/dto"
)

func Invoke(boardDTOs map[string]dto.BoardDTO) domain.PodData {
	podData := domain.PodData{
		Boards: getBoards(boardDTOs),
	}

	return podData
}

func getBoards(boardDTOs map[string]dto.BoardDTO) map[string]*domain.Board {
	boards := make(map[string]*domain.Board)
	for name, boardDTO := range boardDTOs {
		board := &domain.Board{
			Name:    name,
			Packets: boardDTO.GetPackets(),
		}
		boards[board.Name] = board
	}
	return boards
}
