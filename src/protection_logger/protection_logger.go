package protection_logger

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/logger_handler"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

type MessageLogger struct {
	ids           common.Set[string]
	fileName      string
	flushInterval time.Duration
	trace         zerolog.Logger
}

type Config struct {
	FileName      string `toml:"file_name"`
	FlushInterval string `toml:"flush_interval"`
}

// FIXME: instead of "warning" and "fault", use ADE message NUMBER ids as ids.
func NewProtectionLogger(config Config, messagesConfig vehicle.ProtectionConfig) MessageLogger {
	ids := common.NewSet[string]()
	ids.Add(messagesConfig.FaultIdKey)
	ids.Add(messagesConfig.WarningIdKey)

	protectionTrace := trace.With().Str("component", "protectionLogger").Logger()

	flushInterval, err := time.ParseDuration(config.FlushInterval)

	if err != nil {
		protectionTrace.Fatal().Err(err).Msg("error parsing flush interval")
	}

	return MessageLogger{
		ids:           ids,
		fileName:      config.FileName,
		flushInterval: flushInterval,
		trace:         protectionTrace,
	}
}

func (ml *MessageLogger) Ids() common.Set[string] {
	return ml.ids
}

func (ml *MessageLogger) Start(basePath string) chan<- logger_handler.Loggable {
	loggableChan := make(chan logger_handler.Loggable)

	go ml.startLoggingRoutine(loggableChan, basePath)

	return loggableChan
}

func (ml *MessageLogger) startLoggingRoutine(loggableChan <-chan logger_handler.Loggable, basePath string) {
	file := ml.createFile(basePath)
	flushTicker := time.NewTicker(ml.flushInterval)
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
	file, err := logger_handler.NewCSVFile(basePath, ml.fileName)

	if err != nil {
		ml.trace.Fatal().Err(err).Msg("error creating file")
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
