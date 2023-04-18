package data_logger

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/HyperloopUPV-H8/Backend-H8/logger/data_logger/models"
	logger_models "github.com/HyperloopUPV-H8/Backend-H8/logger/models"
	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

type DataLogger struct {
	files map[uint16]map[string]*logger_models.SaveFile

	isRunningMx *sync.Mutex
	isRunning   bool

	config Config
}

func New(config Config) *DataLogger {
	return &DataLogger{
		files: make(map[uint16]map[string]*logger_models.SaveFile),

		isRunningMx: &sync.Mutex{},
		isRunning:   false,

		config: config,
	}
}

func (logger *DataLogger) Start() {
	logger.isRunningMx.Lock()
	defer logger.isRunningMx.Unlock()
	logger.isRunning = true
}

func (logger *DataLogger) Stop() {
	logger.isRunningMx.Lock()
	defer logger.isRunningMx.Unlock()
	logger.isRunning = false
}

func (logger *DataLogger) Update(update vehicle_models.Update) error {
	if !logger.IsRunning() {
		return nil
	}

	return logger.writeValues(update.ID, update.Fields)
}

func (logger *DataLogger) writeValues(id uint16, fields map[string]any) error {
	for name, field := range fields {
		err := logger.writeValue(id, name, field)
		if err != nil {
			return err
		}
	}

	return nil
}

func (logger *DataLogger) writeValue(id uint16, name string, field any) error {
	value := models.NewValue(field)
	file, err := logger.getFile(id, name)
	if err != nil {
		return err
	}

	return file.WriteCSV(value.ToCSV())
}

func (logger *DataLogger) getFile(id uint16, name string) (*logger_models.SaveFile, error) {
	idFiles, ok := logger.files[id]
	if !ok {
		idFiles = make(map[string]*logger_models.SaveFile)
		logger.files[id] = idFiles
	}

	file, ok := idFiles[name]
	if !ok {
		saveFile, err := logger_models.NewSaveFile(logger.getIdPath(id), name+".csv")
		if err != nil {
			return nil, err
		}

		file = saveFile
		idFiles[name] = file
	}

	return file, nil
}

func (logger *DataLogger) getIdPath(id uint16) string {
	return filepath.Join(logger.config.BasePath, fmt.Sprintf("%d", id))
}

func (logger *DataLogger) Flush() (err error) {
	for _, idFiles := range logger.files {
		flushErr := logger_models.FlushFiles(idFiles)
		if flushErr != nil {
			err = flushErr
		}
	}
	return err
}

func (logger *DataLogger) Close() (err error) {
	logger.Stop()

	for _, idFiles := range logger.files {
		closeErr := logger_models.CloseFiles(idFiles)
		if closeErr != nil {
			err = closeErr
		}
	}

	return err
}

func (logger *DataLogger) IsRunning() bool {
	logger.isRunningMx.Lock()
	defer logger.isRunningMx.Unlock()
	return logger.isRunning
}
