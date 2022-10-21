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
		sheetContent, _ := f.GetRows(name)
		tables := getTables(sheetContent)
		sheet := Sheet{
			Name:   name,
			Tables: tables,
		}
		newMap[name] = sheet
	}
	return newMap
}

// POR IMPLEMENTAR
func getTables(sheetContent [][]string) map[string]Table {
	tables := make(map[string]Table)
	var axis [2]int = [2]int{0, 0}
	for axis[0] < len(sheetContent) {
		axis = getInitOfTable(sheetContent, axis)
		tableName := sheetContent[axis[0]][axis[1]]

		var tableContent [][]string //FALTA DE IMPLEMENTAR
		rows := getRowsOfTable(tableContent)

		table := Table{
			Name: tableName,
			Rows: rows,
		}
		tables[tableName] = table

		axis[1]++

	}

	return tables
}

// Otra opción, podía crear nuevo [][]string con lo que queda por recorrer
func getInitOfTable(sheetContent [][]string, axis [2]int) [2]int {
	initString := "[TABLE]"
	finded := false

	for i := axis[0]; i < len(sheetContent); i++ {
		for j := axis[1]; j < len(sheetContent[i]); j++ {
			cellValue := sheetContent[i][j]
			finded = strings.Contains(cellValue, initString)
			if finded {
				axis[0] = i
				axis[1] = j
				break
			}
		}
		if finded {
			break
		}
	}

	return axis
}

// Recibes la hoja cortada
func getRowsOfTable(parcialSheet [][]string) []Row {
	numRows := len(parcialSheet)
	rows := make([]Row, numRows)
	length := len(parcialSheet[0])

	for i := 0; i < length; i++ {
		row := parcialSheet[i]
		rows = append(rows, row)
	}
	return rows
}
