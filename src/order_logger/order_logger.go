package order_logger

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	excel_adapter_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/logger_handler"
)

type OrderLogger struct {
	ids    common.Set[string]
	config Config
}

type Config struct {
	FileName string `toml:"file_name"`
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

	return OrderLogger{
		ids:    ids,
		config: config,
	}
}

func (ol *OrderLogger) Ids() common.Set[string] {
	return ol.ids
}

func (ol *OrderLogger) Start(basePath string) (chan<- logger_handler.Loggable, error) {
	loggableChan := make(chan logger_handler.Loggable)

	go ol.startLoggingRoutine(loggableChan, basePath)

	return loggableChan, nil
}

func (ol *OrderLogger) startLoggingRoutine(loggableChan <-chan logger_handler.Loggable, basePath string) {
	file := ol.createFile(basePath)
	flushTicker := time.NewTicker(time.Second) // PILLAR DE CONF
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
	file, err := logger_handler.NewCSVFile(basePath, ol.config.FileName)

	if err != nil {
		//TODO: trace
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
