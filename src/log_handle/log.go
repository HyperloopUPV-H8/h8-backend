package log_handle

import (
	"fmt"
	golog "log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/log_handle/models"
)

type LogHandle struct {
	source <-chan models.Value
	buffer map[string][]models.Value
	config models.Config
	dumpMx sync.Mutex
	dumps  map[string]*os.File
}

func NewLogger(source <-chan models.Value, config models.Config) *LogHandle {
	return &LogHandle{
		source: source,
		buffer: make(map[string][]models.Value),
		config: config,
		dumpMx: sync.Mutex{},
		dumps:  make(map[string]*os.File),
	}
}

func (log *LogHandle) Run() {
	for {
		if !log.config.Running {
			continue
		}
		select {
		case value := <-log.source:
			log.buffer[value.Name] = append(log.buffer[value.Name], value)
			if len(log.buffer[value.Name]) > int(log.config.DumpSize/log.config.RowSize) {
				log.Dump()
			}
		case <-log.config.Timeout.C:
			log.Dump()
		default:
		}
	}
}

func (log *LogHandle) Start() {
	log.config.Running = true
	log.buffer = make(map[string][]models.Value)
}

func (log *LogHandle) Stop() {
	log.config.Running = false
	log.Dump()
}

func (log *LogHandle) Dump() {
	log.dumpMx.Lock()
	defer log.dumpMx.Unlock()
	for value, buffer := range log.buffer {
		log.writeCSV(value, buffer)
	}
	log.buffer = make(map[string][]models.Value)
}

func (log *LogHandle) writeCSV(value string, buffer []models.Value) {
	file, ok := log.dumps[value]
	if !ok {
		file = log.createFile(value)
		log.dumps[value] = file
	}

	data := ""
	for _, value := range buffer {
		data += fmt.Sprintf("%d,\"%v\"\n", value.Timestamp.Nanosecond(), value.Value)
	}
	file.WriteString(data)
}

func (log *LogHandle) createFile(value string) *os.File {
	os.Mkdir(filepath.Join(log.config.BasePath, value), os.ModeDir)
	file, err := os.Create(filepath.Join(log.config.BasePath, value, strings.ReplaceAll(fmt.Sprintf("%v.csv", time.Now()), " ", "_")))
	if err != nil {
		golog.Fatalf("LogHandle: WriteCSV: %s\n", err)
	}
	return file
}
