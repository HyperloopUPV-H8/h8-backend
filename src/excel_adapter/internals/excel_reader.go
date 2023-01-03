package internals

import (
	"log"
	"strings"

	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/internals/models"
	"github.com/xuri/excelize/v2"
)

const sheetPrefix = "BOARD_ "
const tablePrefix = "[TABLE] "

func GetDocument(file *excelize.File) models.Document {
	sheets := parseSheets(file)
	document := models.Document{
		Sheets: sheets,
	}
	return document
}

func parseSheets(file *excelize.File) map[string]models.Sheet {
	sheets := make(map[string]models.Sheet)
	boards := sheetsFilter(file.GetSheetMap())
	for _, name := range boards {
		cols := getCols(file, name)
		sheets[strings.TrimPrefix(name, sheetPrefix)] = parseSheet(name, cols)
	}
	return sheets

}

func sheetsFilter(sheets map[int]string) map[int]string {
	boards := make(map[int]string)
	for key, name := range sheets {
		if strings.HasPrefix(name, sheetPrefix) {
			boards[key] = name
		}
	}
	return boards
}

func getCols(file *excelize.File, nameSheet string) [][]string {
	cols, err := file.GetCols(nameSheet)
	if err != nil {
		log.Fatalf("error gettings columns: %s\n", err)
	}
	return cols
}

func parseSheet(name string, cols [][]string) models.Sheet {
	tables := make(map[string]models.Table)

	for name, bound := range findTables(cols) {
		println(name)
		tables[name] = parseTable(cols, bound)
	}

	return models.Sheet{
		Tables: tables,
	}
}

func findTables(cols [][]string) map[string][4]int {
	tables := make(map[string][4]int)
	for i, col := range cols {
		for j, cell := range col {
			if strings.HasPrefix(cell, tablePrefix) {
				end := findTableEnd(cols, i, j)
				tables[strings.TrimPrefix(cell, tablePrefix)] = [4]int{i, j, i + end[0], j + end[1] + 2}
			}
		}
	}
	return tables
}

func findTableEnd(cols [][]string, firstCol int, firstRow int) (bound [2]int) {
	widht := findTableWidth(cols, firstCol, firstRow)
	bound[0] = widht

	height := findTableHeight(cols, firstCol, firstRow, widht)
	bound[1] = height
	return
}

func findTableWidth(cols [][]string, firstCol int, firstRow int) int {
	for i, col := range cols[firstCol:] {
		if col[firstRow+1] == "" {
			return i
		} else if i == len(cols[firstCol:])-1 {
			return i + 1
		}
	}
	return -1
}

func findTableHeight(cols [][]string, firstCol int, firstRow int, width int) int {
	maxHeight := 0
	height := 0
	for _, col := range cols[firstCol : width+firstCol] {
		for j, cell := range col[firstRow+2:] { //firstRox+2: first data row
			if cell == "" {
				height = j
				break
			} else if j == len(col[firstRow+2:])-1 {
				height = j + 1
				break
			}
		}

		if height > maxHeight {
			maxHeight = height
		}
	}
	height = maxHeight

	return height
}

func parseTable(cols [][]string, bound [4]int) models.Table {

	rows := make([]models.Row, bound[3]-bound[1]-2)
	for j := 0; j < len(rows); j++ {
		rows[j] = parseRow(cols, j+bound[1]+2, bound[0], bound[2])
	}

	return models.Table{
		Rows: rows,
	}
}

func parseRow(cols [][]string, offset int, start int, end int) models.Row {
	row := make([]string, end-start)
	for i, col := range cols[start:end] {
		if offset < len(col) {
			row[i] = col[offset]
		}
	}
	return row
}
