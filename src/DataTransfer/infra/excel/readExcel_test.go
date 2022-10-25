package excel

import (
	"fmt"
	"log"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestReadExcel(t *testing.T) {
	name := "excelDownloaded.xlsx" //"simple.xlsx" //TO DEFINE
	DocumentExcel := openExcelFile(name)
	mapa := DocumentExcel.GetSheetMap()
	printSheetsName(mapa)
	document := GetDocument(DocumentExcel)
	fmt.Println("Objeto creado: ", document)

}

func openExcelFile(name string) *excelize.File {
	f, err := excelize.OpenFile(name)

	errorsReadingExcel("Opening Excel: ", err)

	return f
}

func errorsReadingExcel(desc string, err error) {
	if err != nil {
		log.Fatal(desc, err)
	}
}

func printSheetsName(sheetNamesMap map[int]string) {

	for index, name := range sheetNamesMap {
		fmt.Println(index, name)
	}
}
