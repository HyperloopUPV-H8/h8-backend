package protection_logger

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/logger_handler"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle"
)

type MessageLogger struct {
	ids    common.Set[string]
	config Config
}

type Config struct {
	FileName string `toml:"file_name"`
}

// FIXME: instead of "warning" and "fault", use ADE message NUMBER ids as ids.
func NewMessageLogger(config Config, messagesConfig vehicle.ProtectionConfig) MessageLogger {
	ids := common.NewSet[string]()
	ids.Add(messagesConfig.FaultIdKey)
	ids.Add(messagesConfig.WarningIdKey)

	return MessageLogger{
		ids:    ids,
		config: config,
	}
}

func (ml *MessageLogger) Ids() common.Set[string] {
	return ml.ids
}

func (ml *MessageLogger) Start(basePath string) (chan<- logger_handler.Loggable, error) {
	loggableChan := make(chan logger_handler.Loggable)

	go ml.startLoggingRoutine(loggableChan, basePath)

	return loggableChan, nil
}

func (ml *MessageLogger) startLoggingRoutine(loggableChan <-chan logger_handler.Loggable, basePath string) {
	file := ml.createFile(basePath)
	flushTicker := time.NewTicker(time.Second) // PILLAR DE CONF
	done := make(chan struct{})
	go ml.startFlushRoutine(flushTicker.C, file, done)

	for loggable := range loggableChan {
		if ml.ids.Has(loggable.Id()) {
			file.Write(loggable.Log())
		}
	}

	done <- struct{}{}
	flushTicker.Stop()
	file.Close()
}

func (ml *MessageLogger) createFile(basePath string) logger_handler.CSVFile {
	file, err := logger_handler.NewCSVFile(basePath, ml.config.FileName)

	if err != nil {
		//TODO: trace
	}

	return file
}

func (ml *MessageLogger) startFlushRoutine(tickerChan <-chan time.Time, file logger_handler.CSVFile, done chan struct{}) {
loop:
	for {
		select {
		case <-tickerChan:
			file.Flush()
		case <-done:
			break loop
		}
	}
}
