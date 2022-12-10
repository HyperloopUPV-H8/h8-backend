package models

import "time"

type Config struct {
	DumpSize uint64
	RowSize  uint64
	Running  bool
	Timeout  time.Ticker
	BasePath string
}
