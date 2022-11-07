package infra

import (
	"path"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/Logger/infra/dto"
)

type Logger struct {
	baseDir    string
	currDir    *Dir
	ticker     *time.Ticker
	EnableChan chan bool
	EntryChan  chan dto.LogPacket
}

func NewLogger(baseDir string, delay time.Duration) Logger {
	return Logger{
		baseDir:    baseDir,
		currDir:    nil,
		ticker:     time.NewTicker(delay),
		EnableChan: make(chan bool),
		EntryChan:  make(chan dto.LogPacket, 100),
	}
}

func (log Logger) Run() {
	go func() {
		for isEnable := range log.EnableChan {
			if isEnable {
				log.record()
			}
		}
	}()
}

func (log Logger) record() {
	log.currDir = NewDir(path.Join(log.baseDir, time.Now().Format("2006-01-02_15-04-05")))
loop:
	for {
		select {
		case <-log.ticker.C:
			log.currDir.Dump()
		case packet := <-log.EntryChan:
			log.addPacket(packet)
		case isEnable := <-log.EnableChan:
			if !isEnable {
				log.currDir.Dump()
				break loop
			}
		}
	}
}

func (log Logger) addPacket(packet dto.LogPacket) {
	for _, measurement := range packet.Values() {
		value := NewValue(packet.Timestamp(), measurement.Data())
		log.currDir.AppendValue(measurement.Name(), value)
	}
}
