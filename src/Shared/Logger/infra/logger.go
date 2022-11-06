package infra

import (
	"path"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain"
)

type Logger struct {
	baseDir    string
	currDir    *Dir
	ticker     *time.Ticker
	EnableChan chan bool
	EntryChan  chan domain.PacketTimestampPair
}

func NewLogger(baseDir string, delay time.Duration) Logger {
	return Logger{
		baseDir:    baseDir,
		currDir:    nil,
		ticker:     time.NewTicker(delay),
		EnableChan: make(chan bool),
		EntryChan:  make(chan domain.PacketTimestampPair, 100),
	}
}

func (log Logger) Run() {
	go func() {
		for {
			isEnable := <-log.EnableChan
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

func (log Logger) addPacket(packet domain.PacketTimestampPair) {
	for name, measurement := range packet.Packet.Measurements {
		value := NewValue(packet.Timestamp, measurement.Value)
		log.currDir.AppendValue(name, value)
	}
}
