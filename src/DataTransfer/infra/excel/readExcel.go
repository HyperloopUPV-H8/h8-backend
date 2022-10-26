package excel

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

func GetDocument(file *excelize.File) Document {
	sheets := getSheets(file)
	document := Document{
		Sheets: sheets,
	}
	return document
}

func getSheets(file *excelize.File) map[string]Sheet {
	newMap := make(map[string]Sheet)
	namesMap := file.GetSheetMap()
	for _, name := range namesMap {
		fmt.Println(name)
		sheet := getSheet(name, file)
		newMap[name] = sheet
	}
	return newMap
}

func getSheet(name string, file *excelize.File) Sheet {
	sheetContent, _ := file.GetRows(name)
	return Sheet{
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

//Deleted func getNumberRows

//Deleted func getCountEmptyCells

func getTables(sheetContent [][]string) map[string]Table {
	tables := make(map[string]Table)
	initCells := getInitOfTables(sheetContent)
	for i := 0; i < len(initCells); i++ {
		table := getTable(sheetContent, initCells[i])
		tables[table.Name] = table
	}

	return tables
}

func getInitOfTables(sheetContent [][]string) [][2]int {
	initString := "[TABLE]"
	var initCells [][2]int

	for i := 0; i < len(sheetContent); i++ {
		for j := 0; j < len(sheetContent[i]); j++ {
			cellValue := sheetContent[i][j]
			if strings.Contains(cellValue, initString) {
				var initCell [2]int = [2]int{i, j}
				initCells = append(initCells, initCell)
			}
		}

	}

	return initCells
}

func getTable(sheetContent [][]string, axis [2]int) Table {
	tableName := getTitleTable(sheetContent[axis[0]][axis[1]])
	headersRow := sheetContent[axis[0]+1]
	rowLength := getMaxRowLength(headersRow, axis[1])
	var initDates [2]int = [2]int{axis[0] + 2, axis[1]}
	rectangularTable := getRectangularTable(sheetContent, rowLength, initDates)
	rows := getRowsOfTable(rectangularTable)

	table := Table{
		Name: tableName,
		Rows: rows,
	}

	return table
}

func getTitleTable(cellContent string) string {
	title := strings.TrimPrefix(cellContent, "[TABLE]")
	return title
}

//Changed getRectangularTable

func getRectangularTable(sheetContent [][]string, rowLength int, initDatesOfTable [2]int) [][]string {
	rectangularTable := make([][]string, 0)

	finalColumn := initDatesOfTable[1] + rowLength - 1 //pos column + number of dates

	for i := initDatesOfTable[0]; i < len(sheetContent); i++ {
		rectangularRow := getRectangularRow(sheetContent[i], initDatesOfTable[1], finalColumn)
		if !isEmpty(rectangularRow) {
			rectangularTable = append(rectangularTable, rectangularRow)
		} else {
			break
		}
	}

	fmt.Println(rectangularTable)
	return rectangularTable
}

func getRectangularRow(fullRow []string, initColumn int, finalColumn int) []string {

	finalIndex := len(fullRow) - 1
	emptySpacesAtEnd := finalColumn - finalIndex
	var row []string

	fmt.Println("finalColumn: ", finalColumn, " finalIndex: ", finalIndex, " emptySpacesAtEnd: ", emptySpacesAtEnd)

	if emptySpacesAtEnd <= 0 {
		row = fullRow[initColumn : finalColumn+1]
	} else {
		row = addEmptyCells(fullRow[initColumn:], emptySpacesAtEnd)
		fmt.Println(len(row))
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

//Deleted func createRectangularTable

// Recibes la hoja cortada
func getRowsOfTable(rectangularTable [][]string) []Row {
	numRows := len(rectangularTable)
	rows := make([]Row, numRows)
	length := len(rectangularTable)

	for i := 0; i < length; i++ {
		rows[i] = rectangularTable[i]
	}
	return rows
}
