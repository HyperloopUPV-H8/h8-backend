package excelRetriever

import (
	"log"
	"path/filepath"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_retriever/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_retriever/infra"
	"github.com/xuri/excelize/v2"
)

func GetExcel(spreadSheetID string, fileName string, filePath string, credentialsPath string) domain.Document {
	infra.FetchExcel(spreadSheetID, fileName, filePath, credentialsPath)
	excel, err := excelize.OpenFile(filepath.Join(filePath, fileName))
	if err != nil {
		log.Fatalf("get excel: got err %s\n", err)
	}
	return infra.GetDocument(excel)
}
