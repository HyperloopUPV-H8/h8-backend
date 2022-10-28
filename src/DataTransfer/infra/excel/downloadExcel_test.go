package excel

import (
	"log"
	"os"
	"testing"
)

func TestDownloadExcel(t *testing.T) {

	// The spreadsheet to request.
	//spreadsheetID := "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms" //El ejemplo
	spreadsheetID := "1nbiLvA0weR_DiLkL9TI90cdLNXlvOAZgikhKIdxbhRk" //Mi spreadsheet con tablas

	filename := "excelDownloaded.xlsx"

	if fileExists(filename) {
		deleteExcel(filename)
	}

	if fileExists(filename) {
		t.Fatalf("file has not been deleted")
	}

	downloadExcel(spreadsheetID, filename)

	if !fileExists(filename) {
		t.Fatalf("file has not been downloaded in %s", filename)
	}
}

func deleteExcel(fileName string) {
	e := os.Remove(fileName)
	if e != nil {
		log.Fatal(e)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
