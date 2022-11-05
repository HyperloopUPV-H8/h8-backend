package infra

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestDownloadExcel(t *testing.T) {
	godotenv.Load("../../../.env")

	// The spreadsheet to request.
	//spreadsheetID := "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms" //El ejemplo
	spreadsheetID := "1nbiLvA0weR_DiLkL9TI90cdLNXlvOAZgikhKIdxbhRk" //Mi spreadsheet con tablas

	fileName := "excelDownloaded.xlsx"

	if fileExists(fileName) {
		deleteExcel(fileName)
	}

	if fileExists(fileName) {
		t.Fatalf("file has not been deleted")
	}

	//downloadExcel(spreadsheetID, filename)
	credentialsPath := "../../../" + os.Getenv("SECRET_FILE_PATH")
	FetchExcel(spreadsheetID, fileName, ".", credentialsPath)

	if !fileExists(fileName) {
		t.Fatalf("file has not been downloaded in %s", fileName)
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
