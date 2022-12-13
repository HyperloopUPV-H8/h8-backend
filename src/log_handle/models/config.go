package models

import "time"

type Config struct {
	DumpSize uint64
	RowSize  uint64
	BasePath string
	Updates  chan map[string]any
	Autosave *time.Ticker
}
