package models

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"sync"

	trace "github.com/rs/zerolog/log"
)

type SaveFile struct {
	fileMx *sync.Mutex
	file   *os.File
	writer *csv.Writer
}

func NewSaveFile(path, name string) (*SaveFile, error) {
	trace.Debug().Str("path", path).Str("name", name).Msg("creating save file")
	if err := os.MkdirAll(path, os.ModeDir); err != nil {
		trace.Error().Stack().Err(err).Str("path", path).Str("name", name).Msg("failed to create directory")
		return nil, err
	}

	file, err := os.Create(filepath.Join(path, name))
	if err != nil {
		trace.Error().Stack().Err(err).Str("path", path).Str("name", name).Msg("failed to create file")
		return nil, err
	}

	return &SaveFile{
		fileMx: &sync.Mutex{},
		file:   file,
		writer: csv.NewWriter(file),
	}, nil
}

func (file *SaveFile) WriteCSV(data []string) error {
	file.fileMx.Lock()
	defer file.fileMx.Unlock()

	trace.Trace().Strs("data", data).Msg("writing csv")
	if err := file.writer.Write(data); err != nil {
		trace.Error().Stack().Err(err).Strs("data", data).Msg("failed to write csv")
		return err
	}
	return nil
}

func (file *SaveFile) Flush() error {
	file.fileMx.Lock()
	defer file.fileMx.Unlock()

	return file.flushUnsafe()
}

func (file *SaveFile) flushUnsafe() error {
	trace.Info().Msg("flushing save file")
	file.writer.Flush()
	return file.writer.Error()
}

func FlushFiles[T comparable](files map[T]*SaveFile) (err error) {
	for _, file := range files {
		flushErr := file.Flush()
		if flushErr != nil {
			err = flushErr
		}
	}
	return err
}

func (file *SaveFile) Close() error {
	file.fileMx.Lock()
	defer file.fileMx.Unlock()

	trace.Info().Msg("closing save file")

	if err := file.flushUnsafe(); err != nil {
		trace.Error().Stack().Err(err).Msg("failed to flush writer")
		return err
	}

	return file.file.Close()
}

func CloseFiles[T comparable](files map[T]*SaveFile) (err error) {
	for _, file := range files {
		closeErr := file.Close()
		if closeErr != nil {
			err = closeErr
		}
	}
	return err
}
