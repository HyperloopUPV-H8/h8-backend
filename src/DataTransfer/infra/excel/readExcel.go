package excel

import (
	"fmt"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ReadExcel() {
	name := "simple.xlsx" //TO DEFINE
	DocumentExcel := OpenExcelFile(name)
	mapa := DocumentExcel.GetSheetMap()
	PrintSheetsName(mapa)
	fmt.Println(mapa)

}

func OpenExcelFile(name string) *excelize.File {
	f, err := excelize.OpenFile(name)

	ErrorsReadingExcel("Opening Excel: ", err)

	return f
}

func ErrorsReadingExcel(desc string, err error) {
	if err != nil {
		log.Fatal(desc, err)
	}
}

func PrintSheetsName(m map[int]string) {

	for index, name := range m {
		fmt.Println(index, name)
	}
}

// Convert Sheets Map To Our Structure
func GetSheets(f *excelize.File) map[string]Sheet {
	newMap := make(map[string]Sheet)
	namesMap := f.GetSheetMap()
	for _, name := range namesMap {
		sheet := getSheet(name, f)
		newMap[name] = sheet
	}
	return newMap
}

func getSheet(name string, f *excelize.File) Sheet {
	sheetContent, _ := f.GetRows(name)
	tables := getTables(sheetContent)
	sheet := Sheet{
		Name:   name,
		Tables: tables,
	}
	return sheet
}

func getMaxRowLength(propertiesRow []string, initColumn int) int {
	maxRowLength := 0
	emptyCell := false
	i := initColumn
	for {
		emptyCell = propertiesRow[i] == ""
		if emptyCell {
			break
		}
		maxRowLength++
		i++
	}

	return maxRowLength
}

func getNumberRows(sheetContent [][]string, initDates [2]int) int {
	numberRows := 0
	emptyCell := false
	i := initDates[0]
	j := initDates[1]
	for {
		emptyCell = sheetContent[i][j] == ""
		if emptyCell {
			break
		}
		numberRows++
		j++
	}

	return numberRows
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
	found := false
	var initCells [][2]int

	for i := 0; i < len(sheetContent); i++ {
		for j := 0; j < len(sheetContent[i]); j++ {
			cellValue := sheetContent[i][j]
			found = strings.Contains(cellValue, initString)
			if found {
				var initCell [2]int = [2]int{i, j}
				initCells = append(initCells, initCell)
				found = false
			}
		}

	}

	return initCells
}

func getTable(sheetContent [][]string, axis [2]int) Table {

	tableName := sheetContent[axis[0]][axis[1]]
	propertiesRow := sheetContent[axis[0]+1]
	rowLengthTable := getMaxRowLength(propertiesRow, axis[1])

	var initDates [2]int = [2]int{axis[0] + 2, axis[1]}
	numberRows := getNumberRows(sheetContent, initDates)
	rectangularTable := getRectangularTable(sheetContent, rowLengthTable, numberRows, initDates)
	rows := getRowsOfTable(rectangularTable)

	table := Table{
		Name: tableName,
		Rows: rows,
	}

	return table
}

// POR IMPLEMENTAR, NO SÉ SI ES CORRECTO
func getRectangularTable(sheetContent [][]string, maxRowLength int, numberRows int, initDatesOfTable [2]int) [][]string {
	rectangularTable := make([][]string, maxRowLength)
	finalRow := initDatesOfTable[0] + numberRows - 1       //pos row + number of rows
	finalCellRow := initDatesOfTable[1] + maxRowLength - 1 //pos column + number of dates

	for i, iR := initDatesOfTable[0], 0; i <= finalRow; i, iR = i+1, iR+1 {
		for j, jR := initDatesOfTable[1], 0; j <= finalCellRow; j, jR = j+1, jR+1 {
			rectangularTable[iR][jR] = sheetContent[i][j] //NO SÉ SI ES CORRECTO
		}
	}

	//check that it is rectangular
	for i := 0; i < len(rectangularTable); i++ {
		for len(rectangularTable[i]) < maxRowLength {
			rectangularTable[i] = append(rectangularTable[i], "")
		}
	}

	return rectangularTable
}

// Recibes la hoja cortada
func getRowsOfTable(rectangularTable [][]string) []Row {
	numRows := len(rectangularTable)
	rows := make([]Row, numRows)
	length := len(rectangularTable)

	for i := 0; i < length; i++ {
		row := rectangularTable[i]
		rows = append(rows, row)
	}
	return rows
}
