package main

import (
	"github.com/HyperloopUPV-H8/Backend-H8/blcu"
	"github.com/HyperloopUPV-H8/Backend-H8/connection_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter"
	"github.com/HyperloopUPV-H8/Backend-H8/file_logger"
	"github.com/HyperloopUPV-H8/Backend-H8/logger_handler"
	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/packet_logger"
	"github.com/HyperloopUPV-H8/Backend-H8/server"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle"
)

type Config struct {
	Excel            excel_adapter.ExcelAdapterConfig
	Connections      connection_transfer.ConnectionTransferConfig
	LoggerHandler    logger_handler.Config `toml:"logger_handler"`
	PacketLogger     packet_logger.Config  `toml:"packet_logger"`
	OrderLogger      file_logger.Config    `toml:"order_logger"`
	ProtectionLogger file_logger.Config    `toml:"protection_logger"`
	Vehicle          vehicle.Config
	DataTransfer     data_transfer.DataTransferConfig `toml:"data_transfer"`
	Orders           struct {
		SendTopic string `toml:"send_topic"`
	}
	Messages message_transfer.MessageTransferConfig
	Server   server.ServerConfig
	BLCU     blcu.BLCUConfig `toml:"blcu"`
}
