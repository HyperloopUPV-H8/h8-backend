package mappers

import (
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/server/domain"
)

func newBoards(boards map[string]excelAdapter.BoardDTO) []domain.Board {
	result := make([]domain.Board, 0, len(boards))
	for _, board := range boards {
		result = append(result, newBoard(board))
	}
	return result
}

func newBoard(board excelAdapter.BoardDTO) domain.Board {
	return domain.Board{
		Name:    board.Name,
		Packets: getPackets(board.GetPackets()),
	}
}
