package infra

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/net/context"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const mimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

func FetchExcel(spreadsheetID string, fileName string, filePath string, credentialsPath string) {
	fileBuf := downloadExcel(spreadsheetID, fileName, credentialsPath)
	saveFileToPath(fileBuf, fileName, filePath)
}

func downloadExcel(spreadsheetID string, fileName string, credentialsPath string) []byte {
	ctx := context.Background()

	driveService, err := drive.NewService(ctx, option.WithCredentialsFile(credentialsPath))
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
	err := os.WriteFile(filepath.Join(filePath, fileName), fileBuf, 0644) //0644 meaning: User: read & write, Group: read, Other: read

	if err != nil {
		log.Fatal("Error saving file to path: ", err)
	}

	return nil
}
