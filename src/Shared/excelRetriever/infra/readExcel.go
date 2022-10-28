package infra

import (
	"strings"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/excelRetriever/domain"
	"github.com/xuri/excelize/v2"
)

func GetDocument(file *excelize.File) domain.Document {
	sheets := getSheets(file)
	document := domain.Document{
		Sheets: sheets,
	}
	return document
}

func getSheets(file *excelize.File) map[string]domain.Sheet {
	newMap := make(map[string]domain.Sheet)
	namesMap := file.GetSheetMap()
	for _, name := range namesMap {
		sheet := getSheet(name, file)
		newMap[name] = sheet
	}
	return newMap
}

func getSheet(name string, file *excelize.File) domain.Sheet {
	sheetContent, _ := file.GetRows(name)
	return domain.Sheet{
		Name:   name,
		Tables: getTables(sheetContent),
	}

}

func getMaxRowLength(propertiesRow []string, initColumn int) int {
	maxRowLength := 0

	for i := initColumn; i < len(propertiesRow); i++ {
		if propertiesRow[i] != "" {
			maxRowLength++
		}
	}
	return maxRowLength
}

func getTables(sheetContent [][]string) map[string]domain.Table {
	tables := make(map[string]domain.Table)
	initCells := getInitOfTables(sheetContent)
	for i := 0; i < len(initCells); i++ {
		table := getTable(sheetContent, initCells[i])
		tables[table.Name] = table
	}

	return tables
}

func getInitOfTables(sheetContent [][]string) [][2]int {

	var initCells [][2]int

	for i := 0; i < len(sheetContent); i++ {

		initCellsOfRow := searchInitsInRow(sheetContent, i)
		initCells = append(initCells, initCellsOfRow...)
	}

	return initCells
}

func searchInitsInRow(sheetContent [][]string, i int) [][2]int {

	initString := "[TABLE] "
	var initCellsOfRow [][2]int

	for j := 0; j < len(sheetContent[i]); j++ {
		cellValue := sheetContent[i][j]
		if strings.Contains(cellValue, initString) {
			var initCell [2]int = [2]int{i, j}
			initCellsOfRow = append(initCellsOfRow, initCell)
		}
	}
	return initCellsOfRow
}

func getTable(sheetContent [][]string, axis [2]int) domain.Table {
	tableName := getTitleTable(sheetContent[axis[0]][axis[1]])
	headersRow := sheetContent[axis[0]+1]
	rowLength := getMaxRowLength(headersRow, axis[1])
	initData := [2]int{axis[0] + 2, axis[1]}
	rectangularTable := getRectangularTable(sheetContent, rowLength, initData)

	table := domain.Table{
		Name: tableName,
		Rows: rectangularTable,
	}

	return table
}

func getTitleTable(cellContent string) string {
	title := strings.TrimPrefix(cellContent, "[TABLE] ")
	return title
}

func getRectangularTable(sheetContent [][]string, rowLength int, initDataOfTable [2]int) [][]string {
	rectangularTable := make([][]string, 0)

	finalColumn := initDataOfTable[1] + rowLength - 1 //pos column + number of data

	for i := initDataOfTable[0]; i < len(sheetContent); i++ {
		rectangularRow := getRectangularRow(sheetContent[i], initDataOfTable[1], finalColumn)
		if !isEmpty(rectangularRow) {
			rectangularTable = append(rectangularTable, rectangularRow)
		} else {
			break
		}
	}

	return rectangularTable
}

func getRectangularRow(fullRow []string, initColumn int, finalColumn int) []string {

	finalIndex := len(fullRow) - 1
	emptySpacesAtEnd := finalColumn - finalIndex
	var row []string

	if emptySpacesAtEnd <= 0 {
		row = fullRow[initColumn : finalColumn+1]
	} else {
		row = addEmptyCells(fullRow[initColumn:], emptySpacesAtEnd)
	}

	return row
}

func isEmpty(row []string) bool {
	countEmpty := 0
	for i := 0; i < len(row); i++ {
		if row[i] == "" {
			countEmpty++
		}
	}
	isEmpty := countEmpty == len(row)

	return isEmpty
}

func addEmptyCells(row []string, emptySpacesAtEnd int) []string {
	for i := 0; i < emptySpacesAtEnd; i++ {
		row = append(row, "")
	}
	return row
}
