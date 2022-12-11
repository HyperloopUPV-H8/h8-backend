package log_handle

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/log_handle/models"
)

type LogHandle struct {
	bufferMx sync.Mutex
	buffer   map[string][]models.Value
	config   models.Config
	filesMx  sync.Mutex
	files    map[string]*os.File
}

func NewLogger(config models.Config) *LogHandle {
	return &LogHandle{
		bufferMx: sync.Mutex{},
		buffer:   make(map[string][]models.Value),
		config:   config,
		filesMx:  sync.Mutex{},
		files:    make(map[string]*os.File),
	}
}

func (logger *LogHandle) Update(updates map[string]any) {
	if !logger.config.IsRunning {
		return
	}

	logger.bufferMx.Lock()
	defer logger.bufferMx.Unlock()
	dump := false
	for name, value := range updates {
		logger.buffer[name] = append(logger.buffer[name], models.Value{
			Value:     value,
			Timestamp: time.Now(),
		})
		if len(logger.buffer[name]) > int(logger.config.DumpSize/logger.config.RowSize) {
			dump = true
		}
	}

	if dump {
		go logger.dump()
	}
}

func (logger *LogHandle) start() {
	logger.config.IsRunning = true
	logger.buffer = make(map[string][]models.Value)
}

func (logger *LogHandle) stop() {
	logger.config.IsRunning = false
	logger.dump()
}

func (logger *LogHandle) dump() {
	logger.filesMx.Lock()
	defer logger.filesMx.Unlock()
	logger.bufferMx.Lock()
	defer logger.bufferMx.Unlock()
	for value, buffer := range logger.buffer {
		logger.writeCSV(value, buffer)
	}
	logger.buffer = make(map[string][]models.Value)
}

func (logger *LogHandle) writeCSV(value string, buffer []models.Value) {
	file := logger.getFile(value)
	data := ""
	for _, value := range buffer {
		data += fmt.Sprintf("%d,\"%v\"\n", value.Timestamp.Nanosecond(), value.Value)
	}
	file.WriteString(data)
}

func (logger *LogHandle) getFile(value string) *os.File {
	if _, ok := logger.files[value]; !ok {
		logger.filesMx.Lock()
		defer logger.filesMx.Unlock()
		logger.files[value] = logger.createFile(value)
	}
	return logger.files[value]
}

func (logger *LogHandle) createFile(value string) *os.File {
	os.Mkdir(filepath.Join(logger.config.BasePath, value), os.ModeDir)
	file, err := os.Create(filepath.Join(logger.config.BasePath, value, strings.ReplaceAll(fmt.Sprintf("%v.csv", time.Now()), " ", "_")))
	if err != nil {
		log.Fatalf("LogHandle: WriteCSV: %s\n", err)
	}
	return file
}

func (logger *LogHandle) HandleRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("log handle: HandleRequest: %s\n", err)
	}

	if string(payload) == "enable" && !logger.config.IsRunning {
		logger.start()
	} else if string(payload) == "disable" && logger.config.IsRunning {
		logger.stop()
	} else if string(payload) != "enable" && string(payload) != "disable" {
		log.Fatalf("log handle: HandleRequest: unexpected body %s\n", payload)
	}
}

func (logger *LogHandle) Close() {
	logger.filesMx.Lock()
	defer logger.filesMx.Unlock()
	for _, file := range logger.files {
		file.Close()
	}
}
