package order_logger

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	excel_adapter_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/logger_handler"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

type OrderLogger struct {
	ids           common.Set[string]
	fileName      string
	flushInterval time.Duration
	trace         zerolog.Logger
}

type Config struct {
	FileName      string `toml:"file_name"`
	FlushInterval string `toml:"flush_interval"`
}

func NewOrderLogger(boards map[string]excel_adapter_models.Board, config Config) OrderLogger {
	ids := common.NewSet[string]()

	for _, board := range boards {
		for _, packet := range board.Packets {
			if packet.Description.Type == "order" {
				ids.Add(packet.Description.ID)
			}
		}
	}

	orderTrace := trace.With().Str("component", "orderLogger").Logger()

	flushInterval, err := time.ParseDuration(config.FlushInterval)

	if err != nil {
		orderTrace.Fatal().Err(err).Msg("error parsing flush interval")
	}

	return OrderLogger{
		ids:           ids,
		fileName:      config.FileName,
		flushInterval: flushInterval,
		trace:         orderTrace,
	}
}

func (ol *OrderLogger) Ids() common.Set[string] {
	return ol.ids
}

func (ol *OrderLogger) Start(basePath string) chan<- logger_handler.Loggable {
	loggableChan := make(chan logger_handler.Loggable)

	go ol.startLoggingRoutine(loggableChan, basePath)

	return loggableChan
}

func (ol *OrderLogger) startLoggingRoutine(loggableChan <-chan logger_handler.Loggable, basePath string) {
	file := ol.createFile(basePath)
	flushTicker := time.NewTicker(ol.flushInterval)
	done := make(chan struct{})
	go ol.startFlushRoutine(flushTicker.C, file, done)

	for loggable := range loggableChan {
		if ol.ids.Has(loggable.Id()) {
			file.Write(loggable.Log())
		}
	}

	done <- struct{}{}
	flushTicker.Stop()
	file.Close()
}

func (ol *OrderLogger) createFile(basePath string) logger_handler.CSVFile {
	file, err := logger_handler.NewCSVFile(basePath, ol.fileName)

	if err != nil {

	}

	return file
}

func (ol *OrderLogger) startFlushRoutine(tickerChan <-chan time.Time, file logger_handler.CSVFile, done chan struct{}) {
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
