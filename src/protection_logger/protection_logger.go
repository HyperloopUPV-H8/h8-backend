package protection_logger

import (
	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/file_logger"
)

func NewProtectionLogger(faultId string, warningId string, config file_logger.Config) file_logger.FileLogger {
	ids := common.NewSet[string]()
	ids.Add(faultId)
	ids.Add(warningId)

	fileLogger := file_logger.NewFileLogger("orderLogger", ids, config)

	return fileLogger
}
