package excelAdapter

import (
	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	excelRetrieverDomain "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_retriever/domain"
)

func GetBoards(structure excelRetrieverDomain.Document) map[string]excelAdapter.BoardDTO {
	boards := make(map[string]excelAdapter.BoardDTO, len(structure.Sheets))
	for name, sheet := range structure.Sheets {
		boards[name] = excelAdapter.NewBoard(sheet)
	}
	return boards
}
