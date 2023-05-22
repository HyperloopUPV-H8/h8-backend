package document

import (
	"github.com/xuri/excelize/v2"
)

func CreateDocument(file *excelize.File) Document {
	fileSheets := getFileSheets(file)

	return Document{
		Sheets: getRectSheets(fileSheets),
	}
}

func getFileSheets(file *excelize.File) map[string][][]string {
	fileSheets := make(map[string][][]string)
	sheetMap := file.GetSheetMap()
	for _, name := range sheetMap {
		sheet, err := file.GetRows(name)

		if err != nil {
			continue
		}

		fileSheets[name] = sheet
	}

	return fileSheets
}

func getRectSheets(sheets map[string][][]string) map[string][][]string {
	rectSheets := make(map[string][][]string)

	for name, sheet := range sheets {
		rectSheets[name] = makeSheetRect(sheet)
	}

	return rectSheets
}

func makeSheetRect(sheet [][]string) [][]string {
	maxLength := 0

	for _, row := range sheet {
		if len(row) > maxLength {
			maxLength = len(row)
		}
	}

	fullRows := make([][]string, 0)
	for _, row := range sheet {
		fullRow := make([]string, maxLength)
		copy(fullRow, row)
		fullRows = append(fullRows, fullRow)
	}

	return fullRows
}
