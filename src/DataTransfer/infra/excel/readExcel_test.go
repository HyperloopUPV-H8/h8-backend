package excel

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

func ReadExcel() {
	name := "simple.xlsx" //TO DEFINE
	DocumentExcel := OpenExcelFile(name)
	mapa := DocumentExcel.GetSheetMap()
	PrintSheetsName(mapa)
	document := GetDocument(DocumentExcel)
	fmt.Println("Objeto creado: ", document)

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

func PrintSheetsName(sheetNamesMap map[int]string) {

	for index, name := range sheetNamesMap {
		fmt.Println(index, name)
	}
}
