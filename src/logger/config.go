package logger

import (
	"github.com/HyperloopUPV-H8/Backend-H8/logger/data_logger"
	"github.com/HyperloopUPV-H8/Backend-H8/logger/message_logger"
)

type Config struct {
	MessageLogger message_logger.Config `toml:"message_logger"`
	DataLogger    data_logger.Config    `toml:"data_logger"`
	FlushInterval string                `toml:"flush_interval"`
	Topics        LoggerTopics          `toml:"topics"`
}

type LoggerTopics struct {
	Enable string `toml:"enable"`
	State  string `toml:"state"`
}
