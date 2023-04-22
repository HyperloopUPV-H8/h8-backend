package packet_logger

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	excel_adapter_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/logger_handler"
)

type PacketLogger struct {
	packetIds    common.Set[string]
	valueIds     common.Set[string]
	valueFilesMx *sync.Mutex
	config       Config
}

type Config struct {
	FolderName     string `toml:"folder_name"`
	PacketFileName string `toml:"packet_file_name"`
}

func NewPacketLogger(boards map[string]excel_adapter_models.Board, config Config) (PacketLogger, error) {
	return PacketLogger{
		packetIds:    getPacketIds(boards),
		valueIds:     getValueIds(boards),
		valueFilesMx: &sync.Mutex{},
		config:       config,
	}, nil
}

func getPacketIds(boards map[string]excel_adapter_models.Board) common.Set[string] {
	ids := common.NewSet[string]()

	for _, board := range boards {
		for _, packet := range board.Packets {
			ids.Add(packet.Description.ID)
		}
	}

	return ids
}

func getValueIds(boards map[string]excel_adapter_models.Board) common.Set[string] {
	ids := common.NewSet[string]()

	for _, board := range boards {
		for _, packet := range board.Packets {
			for _, value := range packet.Values {
				ids.Add(value.ID)
			}
		}
	}

	return ids
}

func (pl *PacketLogger) Ids() common.Set[string] {
	allIds := common.NewSet[string]()

	pl.packetIds.ForEach(func(item string) {
		allIds.Add(item)
	})

	pl.valueIds.ForEach(func(item string) {
		allIds.Add(item)
	})

	return allIds
}

func (pl *PacketLogger) Start(basePath string) (chan<- logger_handler.Loggable, error) {
	loggableChan := make(chan logger_handler.Loggable)

	go pl.startLoggingRoutine(loggableChan, basePath)

	return loggableChan, nil //TODO: change error to something that makes sense
}

func (pl *PacketLogger) startLoggingRoutine(loggableChan <-chan logger_handler.Loggable, basePath string) {
	packetFile := pl.createPacketFile(basePath)
	valueFiles := make(map[string]logger_handler.CSVFile)
	flushTicker := time.NewTicker(time.Second) // PILLAR DE CONF
	done := make(chan struct{})
	go pl.startFlushRoutine(flushTicker.C, packetFile, valueFiles, done)

	for loggable := range loggableChan {
		switch id := loggable.Id(); {
		case pl.packetIds.Has(id):
			packetFile.Write(loggable.Log())
		case pl.valueIds.Has(id):
			pl.valueFilesMx.Lock()
			file := getOrAddFile(valueFiles, filepath.Join(basePath, pl.config.FolderName), id)
			file.Write(loggable.Log())
			pl.valueFilesMx.Unlock()
		default:
			//TODO: trace
		}
	}

	done <- struct{}{}
	flushTicker.Stop()

	closeFiles(packetFile, valueFiles)
}

func (pl *PacketLogger) startFlushRoutine(tickerChan <-chan time.Time, packetFile logger_handler.CSVFile, valueFiles map[string]logger_handler.CSVFile, done chan struct{}) {
loop:
	for {
		select {
		case <-tickerChan:
			packetFile.Flush()
			pl.valueFilesMx.Lock()
			for _, file := range valueFiles {
				file.Flush()
			}
			pl.valueFilesMx.Unlock()
		case <-done:
			break loop
		}
	}
}

func (pl *PacketLogger) createPacketFile(basePath string) logger_handler.CSVFile {
	packetFile, err := logger_handler.NewCSVFile(filepath.Join(basePath, pl.config.FolderName), pl.config.PacketFileName)

	if err != nil {
		//TODO: trace
	}

	return packetFile
}

func getOrAddFile(files map[string]logger_handler.CSVFile, path string, name string) logger_handler.CSVFile {
	file, ok := files[name]
	if !ok {
		newFile, err := logger_handler.NewCSVFile(path, name)

		if err != nil {
			//TODO: TRACE
		}
		files[name] = newFile
		return newFile
	}

	return file
}

func closeFiles(packetFile logger_handler.CSVFile, valueFiles map[string]logger_handler.CSVFile) {
	err := packetFile.Close()

	if err != nil {
		//TODO: trace
	}

	for _, file := range valueFiles {
		err := file.Close()

		if err != nil {
			//TODO: trace
		}
	}
}
