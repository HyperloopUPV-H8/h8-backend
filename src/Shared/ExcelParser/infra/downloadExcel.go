package infra

import (
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/net/context"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const mimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

func DownloadAndSaveExcel(spreadsheetID string, fileName string, filePath string) {
	fileBuf := downloadExcel(spreadsheetID, fileName)
	saveFileToPath(fileBuf, fileName, filePath)
}

func downloadExcel(spreadsheetID string, fileName string) []byte {
	ctx := context.Background()

	driveService, err := drive.NewService(ctx, option.WithCredentialsFile("secret.json"))
	if err != nil {
		log.Fatal("service error: ", err)
	}

	response, err := driveService.Files.Export(spreadsheetID, mimeType).Download()

	if err != nil {
		log.Fatal("http response error: ", err)
	}

	fileBuf, err := io.ReadAll(response.Body)

	if err != nil {
		log.Fatal("Error getting file buffer from httpResponse body: ", err)
	}

	return fileBuf
}

func saveFileToPath(fileBuf []byte, fileName string, filePath string) error {
	err := os.WriteFile(fmt.Sprintf("%v\\%v", fileName, filePath), fileBuf, 0644) //0644 meaning: User: read & write, Group: read, Other: read

	if err != nil {
		log.Fatal("Error saving file to path: ", err)
	}

	return nil
}
