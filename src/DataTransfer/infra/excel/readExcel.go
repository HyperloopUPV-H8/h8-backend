package excel

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

func GetDocument(file *excelize.File) Document {
	sheets := GetSheets(file)
	document := Document{
		Sheets: sheets,
	}
	return document
}

func GetSheets(file *excelize.File) map[string]Sheet {
	newMap := make(map[string]Sheet)
	namesMap := file.GetSheetMap()
	for _, name := range namesMap {
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

func getNumberRows(sheetContent [][]string, initDates [2]int, maxRowLength int) int {
	numberRows := 0
	emptyRow := false

	for i := initDates[0]; i < len(sheetContent); i++ {
		countEmptyCells := getCountEmptyCells(sheetContent, initDates, maxRowLength, i)
		emptyRow = countEmptyCells == maxRowLength
		if emptyRow {
			break
		}
		numberRows++
	}

	fmt.Println("numberRows: ", numberRows)
	return numberRows
}

func getCountEmptyCells(sheetContent [][]string, initDates [2]int, maxRowLength int, row int) int {
	countEmptyCells := 0
	finalJ := initDates[1] + maxRowLength - 1
	emptyCell := false
	for j := initDates[1]; j <= finalJ; j++ {
		emptyCell = sheetContent[row][j] == ""
		if emptyCell {
			countEmptyCells++
		}
	}
	return countEmptyCells
}

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
	numberOfRows := getNumberRows(sheetContent, initDates, rowLength)
	rectangularTable := getRectangularTable(sheetContent, rowLength, numberOfRows, initDates)
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

// devuelve los datos en una tabla rectangular
func getRectangularTable(sheetContent [][]string, maxRowLength int, numberRows int, initDatesOfTable [2]int) [][]string {

	rectangularTable := createRectangularTable(maxRowLength, numberRows)

	finalRow := initDatesOfTable[0] + numberRows - 1       //pos row + number of rows
	finalCellRow := initDatesOfTable[1] + maxRowLength - 1 //pos column + number of dates

	for i, iR := initDatesOfTable[0], 0; i <= finalRow; i, iR = i+1, iR+1 {
		for j, jR := initDatesOfTable[1], 0; j <= finalCellRow; j, jR = j+1, jR+1 {
			if j < len(sheetContent[i]) {
				rectangularTable[iR][jR] = sheetContent[i][j]
			}
		}
	}

	fmt.Println(rectangularTable)
	return rectangularTable
}

func createRectangularTable(maxRowLength int, numberRows int) [][]string {
	rectangularTable := make([][]string, numberRows)

	for i := range rectangularTable {
		rectangularTable[i] = make([]string, maxRowLength)
	}
	return rectangularTable
}

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
