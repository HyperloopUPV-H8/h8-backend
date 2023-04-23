package order_logger

import (
	"github.com/HyperloopUPV-H8/Backend-H8/common"
	excel_adapter_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/file_logger"
)

func NewOrderLogger(boards map[string]excel_adapter_models.Board, config file_logger.Config) file_logger.FileLogger {
	ids := common.NewSet[string]()

	for _, board := range boards {
		for _, packet := range board.Packets {
			if packet.Description.Type == "order" {
				ids.Add(packet.Description.ID)
			}
		}
	}

	fileLogger := file_logger.NewFileLogger("orderLogger", ids, config)

	return fileLogger
}
