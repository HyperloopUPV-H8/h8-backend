package infra

import (
	"fmt"
	golog "log"
	"os"
	"path"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain"
	dataTransfer "github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain/measurement"
)

type Log struct {
	baseDir     string
	files       map[string]*os.File
	writeTicker *time.Ticker
	values      map[string][]ValueTimestamp
	EnableChan  chan bool
	EntryChan   chan domain.PacketTimestampPair
}

func (log Log) Run() {
	go func() {
		for {
			isEnable := <-log.EnableChan
			if isEnable {
				go log.Record()
			}

		}
	}()
}

func NewLog(dir string, measurementNum int, delay time.Duration) *Log {
	return &Log{
		baseDir:     path.Join(dir, time.Now().Format("2006/01/02-15:04:05")),
		files:       make(map[string]*os.File, measurementNum),
		writeTicker: time.NewTicker(delay),
		values:      make(map[string][]ValueTimestamp, measurementNum),
		EnableChan:  make(chan bool, 1),
		EntryChan:   make(chan domain.PacketTimestampPair, 100),
	}
}

func (log *Log) AddValues(packet dataTransfer.PacketTimestampPair) {
	for _, measurement := range packet.Packet.Measurements {
		log.addValue(packet.Timestamp, measurement)
	}
}

func (log *Log) addValue(timestamp time.Time, value measurement.Measurement) {
	values, exists := log.values[value.Name]

	if !exists {
		values = make([]ValueTimestamp, 100)
	}

	log.values[value.Name] = append(values, NewValue(timestamp, value.Value))
}

func (log *Log) ToString(valueName string) (data string) {
	for _, value := range log.values[valueName] {
		data += value.ToString() + "\n"
	}
	return data
}

func (log *Log) Write() {
	for name := range log.values {
		fmt.Println(name)
		log.writeValues(name)
	}
	log.cleanValues()
}

func (log *Log) writeValues(valueName string) {
	file, exists := log.files[valueName]
	var err error
	if !exists {
		file, err = os.Create(path.Join(log.baseDir, valueName))
		fmt.Println(err)
	}

	if err != nil {
		golog.Fatalf("write values: %s\n", err)
	}

	file.WriteString(log.ToString(valueName))
}

func (log *Log) cleanValues() {
	log.values = make(map[string][]ValueTimestamp, len(log.values))
}

func (log *Log) Close() {
	for _, file := range log.files {
		file.Close()
	}
	log.writeTicker.Stop()
}

func (log Log) Record() {
loop:
	for {
		select {
		case <-log.writeTicker.C:
			log.Write()
		case entries := <-log.EntryChan:
			log.AddValues(entries)
		case isEnabled := <-log.EnableChan:
			if !isEnabled {
				log.Write()
				log.Close()
				break loop
			}
		}
	}
}
