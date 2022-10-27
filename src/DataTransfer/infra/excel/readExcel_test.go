package excel

import (
	"fmt"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestReadExcel(t *testing.T) {
	name := "excelDownloaded.xlsx" //TO DEFINE
	DocumentExcel, err := openExcelFile(name)

	if err != nil {
		t.Fatalf("Couldn't open the file")
	}

	document := GetDocument(DocumentExcel)
	fmt.Println("Objeto creado: ", document)
	//correctObject := {map[Hoja 1:{Hoja 1 map[ NOMBRE1:{ NOMBRE1 [[1 1 1 1] [2 2 2 2] [3 3 3 3]]}]} Hoja 2:{Hoja 2 map[ NOMBRE1:{ NOMBRE1 [[1 1 1] [2 2 ] [3 3 3]]}  NOMBRE2:{ NOMBRE2 [[4 4 4] [5 5 5] [6 6 6]]}]} Hoja 3:{Hoja 3 map[ NOMBRE1:{ NOMBRE1 [[ 1 1 1] [2 2 2 2] [3 3  3]]}  NOMBRE2:{ NOMBRE2 [[4 4 4 ] [5 5 5 5] [6 6 6 ]]}]}]}

}

func openExcelFile(name string) (*excelize.File, error) {
	f, err := excelize.OpenFile(name)

	return f, err
}
