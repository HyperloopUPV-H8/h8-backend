package podDataCreator

import (
	excelRetreiver "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/excelRetreiver"
	domain "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/podDataCreator/infra/excelAdapter"
)

func New() domain.PodData {
	structure := excelRetreiver.GetStructure()
	podData := domain.PodData{}
	podData.Boards = getBoards(structure)

	return podData
}

func getBoards(structure excelRetreiver.Structure) map[string]*domain.Board {
	boards := make(map[string]*domain.Board)
	for name, sheet := range structure.Sheets {
		board := &domain.Board{
			Name:    name,
			Packets: excelAdapter.GetPackets(sheet.Tables),
		}
		boards[board.Name] = board
	}
	return boards
}
