package data_logger

import (
	"fmt"
	"path/filepath"

	"github.com/HyperloopUPV-H8/Backend-H8/logger/data/internals"
	logger_models "github.com/HyperloopUPV-H8/Backend-H8/logger/models"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/data"
)

type DataLogger struct {
	file *logger_models.SaveFile
	dirs map[uint16]*internals.PacketDir

	isRunning bool

	config Config
}

func New(config Config) (*DataLogger, error) {
	file, err := logger_models.NewSaveFile(config.BasePath, config.FileName)
	if err != nil {
		return nil, err
	}

	return &DataLogger{
		file: file,
		dirs: make(map[uint16]*internals.PacketDir),

		isRunning: false,

		config: config,
	}, nil
}

func (logger *DataLogger) Start() error {
	logger.isRunning = true
	return nil
}

func (logger *DataLogger) Stop() error {
	logger.isRunning = false
	return nil
}

func (logger *DataLogger) Update(dta packet.Packet) error {
	if !logger.isRunning {
		return nil
	}

	payload, ok := dta.Payload.(data.Payload)
	if !ok {
		return fmt.Errorf("invalid payload type")
	}

	return logger.writePacket(dta.Metadata, payload)
}

func (logger *DataLogger) writePacket(meta packet.Metadata, payload data.Payload) error {
	err := logger.writeGlobal(meta, payload)
	if err != nil {
		return err
	}

	err = logger.writeValues(meta, payload)
	if err != nil {
		return err
	}

	return nil
}

func (logger *DataLogger) writeGlobal(meta packet.Metadata, payload data.Payload) error {
	record := logger.getGlobalRecord(meta, payload)
	return logger.file.WriteCSV(record)
}

func (logger *DataLogger) getGlobalRecord(meta packet.Metadata, payload data.Payload) []string {
	return []string{
		fmt.Sprint(meta.Timestamp.UnixNano()),
		meta.From,
		meta.To,
		fmt.Sprint(meta.SeqNum),
		fmt.Sprint(meta.ID),
		fmt.Sprint(payload.Values),
	}
}

func (logger *DataLogger) writeValues(meta packet.Metadata, payload data.Payload) error {
	dir := logger.getDir(meta.ID)
	return dir.Write(meta, payload)
}

func (logger *DataLogger) getDir(id uint16) *internals.PacketDir {
	dir, ok := logger.dirs[id]
	if !ok {
		path := filepath.Join(logger.config.BasePath, fmt.Sprint(id))
		dir = internals.NewDir(path)
		logger.dirs[id] = dir
	}

	return dir
}

func (logger *DataLogger) Flush() (err error) {
	for _, dir := range logger.dirs {
		flushErr := dir.Flush()
		if flushErr != nil {
			err = flushErr
		}
	}

	flushErr := logger.file.Flush()
	if flushErr != nil {
		err = flushErr
	}

	return err
}

func (logger *DataLogger) Close() (err error) {
	logger.Stop()

	for _, dir := range logger.dirs {
		closeErr := dir.Close()
		if closeErr != nil {
			err = closeErr
		}
	}

	closeErr := logger.file.Close()
	if closeErr != nil {
		err = closeErr
	}

	return err
}
