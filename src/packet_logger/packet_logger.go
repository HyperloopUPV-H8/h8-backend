package packet_logger

import (
	"github.com/HyperloopUPV-H8/Backend-H8/common"
	excel_adapter_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/file_logger"
)

func NewPacketLogger(boards map[string]excel_adapter_models.Board, config file_logger.Config) file_logger.FileLogger {
	ids := getIds(boards)

	return file_logger.NewFileLogger("packetLogger", ids, config)
}

func getIds(boards map[string]excel_adapter_models.Board) common.Set[string] {
	ids := common.NewSet[string]()

	for _, board := range boards {
		for _, packet := range board.Packets {
			ids.Add(packet.Description.ID)
		}
	}

	return ids
}
