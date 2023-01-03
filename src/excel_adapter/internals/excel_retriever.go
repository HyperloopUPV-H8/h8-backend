package internals

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const sheetsMimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

func DownloadFile(id string, path string, name string) {
	saveFile(getFile(getClient(), id, sheetsMimeType), path, name)
}

func getClient() *drive.Service {
	ctx := context.Background()

	client, err := drive.NewService(ctx, option.WithCredentialsFile(os.Getenv("CREDENTIALS")))
	if err != nil {
		log.Fatalf("excel retriever: getClient: %s\n", err)
	}

	return client
}

func getFile(client *drive.Service, id string, mimeType string) []byte {
	resp, err := client.Files.Export(id, mimeType).Download()
	if err != nil {
		log.Fatalf("excel retriever: getFile: %s\n", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("excel retriever: getFile: %s\n", err)
	}

	return data
}

func saveFile(content []byte, path string, name string) {
	err := os.Mkdir(path, os.ModeDir)
	if !os.IsExist(err) {
		log.Fatalf("excel retriever: saveFile: %s\n", err)
	}
	err = os.WriteFile(filepath.Join(path, name), content, 0644) // rw-r--r--
	if err != nil {
		log.Fatalf("excel retriever: saveFile: %s\n", err)
	}
}
