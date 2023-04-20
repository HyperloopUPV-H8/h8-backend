package message_logger

import (
	"fmt"

	logger_models "github.com/HyperloopUPV-H8/Backend-H8/logger/models"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/message"
)

type MessageLogger struct {
	file *logger_models.SaveFile

	isRunning bool

	config Config
}

func New(config Config) (*MessageLogger, error) {
	file, err := logger_models.NewSaveFile(config.BasePath, config.FileName)
	if err != nil {
		return nil, err
	}

	logger := &MessageLogger{
		file: file,

		config: config,
	}

	return logger, logger.writeHeader()
}

func (logger *MessageLogger) writeHeader() error {
	return logger.writeMsg([]string{
		"timestamp",
		"source",
		"destination",
		"seq_num",
		"id",
		"message",
	})
}

func (logger *MessageLogger) Start() error {
	logger.isRunning = true
	return nil
}

func (logger *MessageLogger) Stop() error {
	logger.isRunning = false
	return nil
}

func (logger *MessageLogger) Update(msg packet.Packet) error {
	if !logger.isRunning {
		return nil
	}

	payload, ok := msg.Payload.(message.Payload)
	if !ok {
		return fmt.Errorf("payload is not a message payload")
	}

	record := logger.toCSV(msg.Metadata, payload)

	return logger.writeMsg(record)
}

func (logger *MessageLogger) toCSV(meta packet.Metadata, payload message.Payload) []string {
	return []string{
		fmt.Sprint(meta.Timestamp.UnixNano()),
		meta.From,
		meta.To,
		fmt.Sprint(meta.SeqNum),
		fmt.Sprint(meta.ID),
		payload.Data.String(),
	}
}

func (logger *MessageLogger) writeMsg(msg []string) error {
	return logger.file.WriteCSV(msg)
}

func (logger *MessageLogger) Close() error {
	return logger.file.Close()
}

func (logger *MessageLogger) Flush() error {
	return logger.file.Flush()
}
