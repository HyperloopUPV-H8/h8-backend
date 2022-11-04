package infra

import (
	"fmt"
	"log"
	"strings"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_retriever/domain"
	"github.com/xuri/excelize/v2"
)

const sheetPrefix = "BOARD "
const tablePrefix = "[TABLE] "

func GetDocument(file *excelize.File) domain.Document {
	sheets := ParseSheets(file)
	document := domain.Document{
		Sheets: sheets,
	}
	return document
}

func ParseSheets(file *excelize.File) map[string]domain.Sheet {
	sheets := make(map[string]domain.Sheet)
	boards := sheetsFilter(file.GetSheetMap())
	fmt.Println(boards)
	for _, name := range boards {
		cols, err := file.GetCols(name)
		fmt.Println(cols)
		if err != nil {
			log.Fatalf("get rows: %s\n", err)
		}
		fmt.Println(name)
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

func parseSheet(name string, cols [][]string) domain.Sheet {
	fmt.Println("Hola")
	tables := make(map[string]domain.Table)
	for name, bound := range findTables(cols) {
		tables[name] = parseTable(name, cols, bound)
	}

	return domain.Sheet{
		Name:   strings.TrimPrefix(name, sheetPrefix),
		Tables: tables,
	}
}

func findTables(cols [][]string) map[string][4]int {
	tables := make(map[string][4]int)
	for i, col := range cols {
		for j, cell := range col {
			if strings.HasPrefix(cell, tablePrefix) {
				end := findTableEnd(cols, i, j)
				tables[strings.TrimPrefix(cell, tablePrefix)] = [4]int{i, j, i + end[0], j + end[1]}
			}
		}
	}
	return tables
}

func findTableEnd(cols [][]string, firstCol int, firstRow int) (bound [2]int) {
	for i, col := range cols[firstCol:] {
		if col[firstRow+1] == "" {
			bound[0] = i
			break
		}
	}

	for j, cell := range cols[firstCol][firstRow:] {
		if cell == "" {
			bound[1] = j
			break
		}
	}

	return
}

func parseTable(name string, cols [][]string, bound [4]int) domain.Table {
	fmt.Println(bound[3] - bound[1] - 2)
	rows := make([]domain.Row, bound[3]-bound[1]-2)
	for j := 0; j < len(rows); j++ {
		rows[j] = parseRow(cols, j+bound[1]+2, bound[0], bound[2])
	}

	return domain.Table{
		Name: name,
		Rows: rows,
	}
}

func parseRow(cols [][]string, offset int, start int, end int) domain.Row {
	row := make([]string, end-start)
	for i, col := range cols[start:end] {
		row[i] = col[offset]
	}
	return row
}
