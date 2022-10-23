package excel

import (
	"fmt"
	"log"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/infra/excel"
	"github.com/xuri/excelize/v2"
)

func ReadExcel() {
	name := "simple.xlsx" //TO DEFINE
	DocumentExcel := OpenExcelFile(name)
	mapa := DocumentExcel.GetSheetMap()
	PrintSheetsName(mapa)
	//fmt.Println(mapa)
	//hoja1, _ := DocumentExcel.GetRows("Hoja1")
	//fmt.Println(hoja1)
	document := excel.GetDocument(DocumentExcel)
	fmt.Println(document)

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
