package main

import (
	"github.com/HyperloopUPV-H8/Backend-H8/board/boards/blcu"
	"github.com/HyperloopUPV-H8/Backend-H8/connection_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/excel_adapter"
	"github.com/HyperloopUPV-H8/Backend-H8/logger"
	"github.com/HyperloopUPV-H8/Backend-H8/message_transfer"
	"github.com/HyperloopUPV-H8/Backend-H8/server"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle"
)

type Config struct {
	Excel        excel_adapter.ExcelAdapterConfig
	Connections  connection_transfer.ConnectionTransferConfig
	Logger       logger.LoggerConfig
	Vehicle      vehicle.VehicleConfig
	DataTransfer data_transfer.DataTransferConfig `toml:"data_transfer"`
	Orders       struct {
		SendTopic string `toml:"send_topic"`
	}
	Messages message_transfer.MessageTransferConfig
	Server   server.ServerConfig
	BLCU     blcu.BLCUConfig `toml:"blcu"`
}
