package infra

import (
	"reflect"
	"testing"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/ExcelParser/domain/document"
	"github.com/xuri/excelize/v2"
)

func TestReadExcel(t *testing.T) {
	name := "excelDownloaded.xlsx" //TO DEFINE
	DocumentExcel, err := openExcelFile(name)

	if err != nil {
		t.Fatalf("couldn't open the file")
	}

	document := GetDocument(DocumentExcel)

	correctObject := getCorrectDocument()

	areEquals := reflect.DeepEqual(document, correctObject)

	if !areEquals {
		t.Fatalf("objects are not equal")
	}
}

func openExcelFile(name string) (*excelize.File, error) {
	f, err := excelize.OpenFile(name)

	return f, err
}

func getCorrectDocument() document.Document {

	correctRows1_1 := [][]string{{"1", "1", "1", "1"}, {"2", "2", "2", "2"}, {"3", "3", "3", "3"}}

	correctTable1_1 := document.Table{
		Name: "NOMBRE1",
		Rows: correctRows1_1,
	}
	correctSheet1 := document.Sheet{
		Name: "Hoja 1",
		Tables: map[string]document.Table{
			"NOMBRE1": correctTable1_1,
		},
	}
	correctRows2_1 := [][]string{{"1", "1", "1"}, {"2", "2", ""}, {"3", "3", "3"}}

	correctTable2_1 := document.Table{
		Name: "NOMBRE1",
		Rows: correctRows2_1,
	}

	correctRows2_2 := [][]string{{"4", "4", "4"}, {"5", "5", "5"}, {"6", "6", "6"}}

	correctTable2_2 := document.Table{
		Name: "NOMBRE2",
		Rows: correctRows2_2,
	}

	correctSheet2 := document.Sheet{
		Name: "Hoja 2",
		Tables: map[string]document.Table{
			"NOMBRE1": correctTable2_1,
			"NOMBRE2": correctTable2_2,
		},
	}

	correctRows3_1 := [][]string{{"", "1", "1", "1"}, {"2", "2", "2", "2"}, {"3", "3", "", "3"}}

	correctTable3_1 := document.Table{
		Name: "NOMBRE1",
		Rows: correctRows3_1,
	}

	correctRows3_2 := [][]string{{"4", "4", "4", ""}, {"5", "5", "5", "5"}, {"6", "6", "6", ""}}

	correctTable3_2 := document.Table{
		Name: "NOMBRE2",
		Rows: correctRows3_2,
	}

	correctSheet3 := document.Sheet{
		Name: "Hoja 3",
		Tables: map[string]document.Table{
			"NOMBRE1": correctTable3_1,
			"NOMBRE2": correctTable3_2,
		},
	}

	document := document.Document{
		Sheets: map[string]document.Sheet{
			"Hoja 1": correctSheet1,
			"Hoja 2": correctSheet2,
			"Hoja 3": correctSheet3,
		},
	}

	return document
}
