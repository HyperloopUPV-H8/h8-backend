package application

import (
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/application/interfaces"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain/board"
)

func GetBoards(structure domain.Document) map[string]interfaces.Board {
	boards := make(map[string]interfaces.Board, len(structure.Sheets))
	for name, sheet := range structure.Sheets {
		boards[name] = board.NewBoard(sheet)
	}
	return boards
}
