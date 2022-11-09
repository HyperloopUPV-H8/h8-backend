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
	ValueChan  chan dto.LogValue
}

func NewLogger(baseDir string, delay time.Duration) Logger {
	return Logger{
		baseDir:    baseDir,
		currDir:    nil,
		ticker:     time.NewTicker(delay),
		EnableChan: make(chan bool),
		ValueChan:  make(chan dto.LogValue, 100),
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
		case value := <-log.ValueChan:
			log.addValue(value)
		case isEnable := <-log.EnableChan:
			if !isEnable {
				log.currDir.Dump()
				break loop
			}
		}
	}
}

func (log Logger) addValue(value dto.LogValue) {
	val := NewValue(value.Timestamp(), value.Data())
	log.currDir.AppendValue(value.Name(), val)
}
