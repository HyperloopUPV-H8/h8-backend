package application

import (
	"log"
	"os"
	"path/filepath"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain/board"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain/document"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/infra"
	"github.com/xuri/excelize/v2"
)

func GetBoards(structure document.Document) map[string]board.Board {
	boards := make(map[string]board.Board, len(structure.Sheets))
	for name, sheet := range structure.Sheets {
		boards[name] = board.New(sheet)
	}
	return boards
}

func GetExcel(fileName string, filePath string) document.Document {
	infra.FetchExcel(os.Getenv("SPREADSHEET_ID"), fileName, filePath)
	excel, err := excelize.OpenFile(filepath.Join(filePath, fileName))
	if err != nil {
		log.Fatalf("get excel: got err %s\n", err)
	}
	return infra.GetDocument(excel)
}
