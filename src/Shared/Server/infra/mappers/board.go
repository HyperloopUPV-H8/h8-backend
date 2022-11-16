package mappers

import (
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/server/domain"
)

func newBoards(boards map[string]excelAdapter.BoardDTO) map[string]domain.Board {
	result := make(map[string]domain.Board, len(boards))
	for name, board := range boards {
		result[name] = newBoard(board)
	}
	return result
}

func newBoard(board excelAdapter.BoardDTO) domain.Board {
	return domain.Board{
		Name:    board.Name,
		Packets: getPackets(board.GetPackets()),
	}
}
