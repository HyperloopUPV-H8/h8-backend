package message_logger

import (
	"github.com/HyperloopUPV-H8/Backend-H8/logger/message_logger/models"
	logger_models "github.com/HyperloopUPV-H8/Backend-H8/logger/models"
)

type MessageLogger struct {
	file *logger_models.SaveFile

	config Config
}

func New(config Config) (*MessageLogger, error) {
	file, err := logger_models.NewSaveFile(config.BasePath, config.FileName)
	if err != nil {
		return nil, err
	}

	return &MessageLogger{
		file: file,

		config: config,
	}, nil
}

func (logger *MessageLogger) Update(message string) error {
	msg := models.NewMessage(message)

	return logger.writeMsg(msg)
}

func (logger *MessageLogger) writeMsg(msg models.Message) error {
	return logger.file.WriteCSV(msg.ToCSV())
}

func (logger *MessageLogger) Close() error {
	return logger.file.Close()
}

func (logger *MessageLogger) Flush() error {
	return logger.file.Flush()
}
