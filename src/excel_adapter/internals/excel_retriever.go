package internals

import (
	"context"
	"io"
	"os"
	"path/filepath"

	trace "github.com/rs/zerolog/log"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const SHEETS_MIME_TYPE = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

func DownloadFile(id string, path string, name string) error {
	trace.Trace().Str("id", id).Str("path", path).Str("name", name).Msg("download file")
	client, errClient := getClient()
	if errClient != nil {
		trace.Error().Str("id", id).Str("path", path).Str("name", name).Stack().Err(errClient).Msg("")
		return errClient
	}

	file, errFile := getFile(client, id, SHEETS_MIME_TYPE)
	if errFile != nil {
		trace.Error().Str("id", id).Str("path", path).Str("name", name).Stack().Err(errFile).Msg("")
		return errFile
	}

	errSaving := saveFile(file, path, name)
	if errSaving != nil {
		trace.Error().Str("id", id).Str("path", path).Str("name", name).Stack().Err(errSaving).Msg("")
	}

	return errSaving
}

func getClient() (*drive.Service, error) {
	trace.Trace().Msg("get client")
	ctx := context.Background()

	client, err := drive.NewService(ctx, option.WithCredentialsFile(os.Getenv("EXCEL_ADAPTER_CREDENTIALS")))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getFile(client *drive.Service, id string, mimeType string) ([]byte, error) {
	trace.Trace().Str("id", id).Msg("get file")
	resp, err := client.Files.Export(id, mimeType).Download()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func saveFile(content []byte, path string, name string) error {
	trace.Trace().Str("path", path).Str("name", name).Msg("save file")

	err := os.Mkdir(path, os.ModeDir)
	if !os.IsExist(err) {
		return err
	}

	err = os.WriteFile(filepath.Join(path, name), content, 0644) // rw-r--r--
	if err != nil {
		return err
	}
	return nil
}
