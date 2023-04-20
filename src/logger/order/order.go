package message_logger

import (
	"fmt"

	logger_models "github.com/HyperloopUPV-H8/Backend-H8/logger/models"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/HyperloopUPV-H8/Backend-H8/packet/order"
)

type OrderLogger struct {
	file *logger_models.SaveFile

	isRunning bool

	config Config
}

func New(config Config) (*OrderLogger, error) {
	file, err := logger_models.NewSaveFile(config.BasePath, config.FileName)
	if err != nil {
		return nil, err
	}

	logger := &OrderLogger{
		file: file,

		config: config,
	}

	return logger, logger.writeHeader()
}

func (logger *OrderLogger) writeHeader() error {
	return logger.writeMsg([]string{
		"timestamp",
		"source",
		"destination",
		"seq_num",
		"id",
		"message",
	})
}

func (logger *OrderLogger) Start() error {
	logger.isRunning = true
	return nil
}

func (logger *OrderLogger) Stop() error {
	logger.isRunning = false
	return nil
}

func (logger *OrderLogger) Update(msg packet.Packet) error {
	if !logger.isRunning {
		return nil
	}

	payload, ok := msg.Payload.(order.Payload)
	if !ok {
		return fmt.Errorf("payload is not a message payload")
	}

	record := logger.toCSV(msg.Metadata, payload)

	return logger.writeMsg(record)
}

func (logger *OrderLogger) toCSV(meta packet.Metadata, payload order.Payload) []string {
	return []string{
		fmt.Sprint(meta.Timestamp.UnixNano()),
		meta.From,
		meta.To,
		fmt.Sprint(meta.SeqNum),
		fmt.Sprint(meta.ID),
		fmt.Sprint(payload.Values),
		fmt.Sprint(payload.Enabled),
	}
}

func (logger *OrderLogger) writeMsg(msg []string) error {
	return logger.file.WriteCSV(msg)
}

func (logger *OrderLogger) Close() error {
	return logger.file.Close()
}

func (logger *OrderLogger) Flush() error {
	return logger.file.Flush()
}
