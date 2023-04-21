package packet_logger

import (
	"errors"
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	excel_adapter_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/logger"
)

type PacketLogger struct {
	packetIds common.Set[string]
	valueIds common.Set[string]
	packetFile logger.SaveFile
	valueFiles map[string]logger.SaveFile
	isRunning  bool
	config     Config
}

type Config struct {
	BasePath string `toml:"base_path"`
}

func NewPacketLogger(boards map[string]excel_adapter_models.Board, config Config) (PacketLogger, error) {
	packetFile, err := logger.NewSaveFile(config.BasePath, config.FileName)

	if err != nil {
		return PacketLogger{}, err
	}

	return PacketLogger{
		packetIds: getPacketIds(boards),
		valueIds: getValueIds(boards),

		packetFile: packetFile,
		// FIXME: Es posible que la SaveFile tenga que ser
		// un puntero
		valueFiles: make(map[string]logger.SaveFile),

		isRunning: false,

		config: config,
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

func (pl *PacketLogger) Start() error {
	pl.isRunning = true
	return nil
}

func (pl *PacketLogger) Stop() error {
	pl.isRunning = false
	return nil
}

func (pl *PacketLogger) Log(loggable logger.Loggable) error {
	if !pl.isRunning {
		return nil
	}

	if pl.packetIds.Has(loggable.Id()) {
		pl.packetFile.WriteCSV(loggable.Log())
	} else if pl.valueIds.Has(loggable.Id()) {
		file, ok := pl.valueFiles[loggable.Id()]

		if !ok {
			file = logger.NewSaveFile(...)
			return 
		}

		file.WriteCSV(loggable.Log())
	} else {
		return fmt.Errorf("loggable id not recognized: %+v", loggable)
	}

}
