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

func DownloadFile(id string, path string, name string) error {
	client, errClient := getClient()
	if errClient != nil {
		return errClient
	}

	file, errFile := getFile(client, id, sheetsMimeType)
	if errFile != nil {
		return errFile
	}

	errSaving := saveFile(file, path, name)

	return errSaving
}

func getClient() (*drive.Service, error) {
	ctx := context.Background()

	client, err := drive.NewService(ctx, option.WithCredentialsFile(os.Getenv("CREDENTIALS")))
	if err != nil {
		log.Printf("excel retriever: the client could not be obtained")
		return nil, err
	}

	return client, nil
}

func getFile(client *drive.Service, id string, mimeType string) ([]byte, error) {
	resp, err := client.Files.Export(id, mimeType).Download()
	if err != nil {
		log.Printf("excel retriever: getFile: could not download the file: %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("excel retriever: getFile: could not read the http response body")
		return nil, err
	}

	return data, nil
}

func saveFile(content []byte, path string, name string) error {
	err := os.Mkdir(path, os.ModeDir)
	if !os.IsExist(err) {
		log.Printf("excel retriever: saveFile: could not create the directory")
		return err
	}
	err = os.WriteFile(filepath.Join(path, name), content, 0644) // rw-r--r--
	if err != nil {
		log.Printf("excel retriever: saveFile: could not write the file in the directory")
		return err
	}
	return nil
}
