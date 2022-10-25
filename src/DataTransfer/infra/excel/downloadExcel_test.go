package excel

import "testing"

func TestDownloadExcel(t *testing.T) {
	// The spreadsheet to request.
	//spreadsheetID := "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms" //El ejemplo
	//spreadsheetID := "1lVcaXeInAThCDvsoV69MVlrd0TzsAi2209v0jvgPbeo" //Mi spreadsheet de prueba
	spreadsheetID := "1nbiLvA0weR_DiLkL9TI90cdLNXlvOAZgikhKIdxbhRk" //Mi spreadsheet con tablas

	filename := "excelDownloaded.xlsx"
	downloadExcel(spreadsheetID, filename)
}
