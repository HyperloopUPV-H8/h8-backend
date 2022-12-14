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
	flushing bool
	buffer   map[string][]models.Value
	config   models.Config
	fileMx   sync.Mutex
	files    map[string]*os.File
	done     chan struct{}
	running  bool
}

func NewLogger(config models.Config) *LogHandle {
	return &LogHandle{
		flushing: false,
		buffer:   make(map[string][]models.Value),
		config:   config,
		fileMx:   sync.Mutex{},
		files:    make(map[string]*os.File),
		done:     make(chan struct{}),
	}
}

func (logger *LogHandle) run() {
	logger.running = true
	defer func() { logger.running = false }()
	for {
		select {
		case update := <-logger.config.Updates:
			for name, value := range update {
				logger.buffer[name] = append(logger.buffer[name], models.Value{
					Value:     value,
					Timestamp: time.Now(),
				})
			}

			logger.checkDump()
		case <-logger.config.Autosave.C:
			logger.flush()
		case <-logger.done:
			return
		}
	}
}

func (logger *LogHandle) checkDump() {
	for _, buf := range logger.buffer {
		if len(buf) > int(logger.config.DumpSize/logger.config.RowSize) {
			logger.flush()
			break
		}
	}
}

func (logger *LogHandle) Update(values map[string]any) {
	if logger.running {
		logger.config.Updates <- values
	}
}

func (logger *LogHandle) start() {
	logger.buffer = make(map[string][]models.Value)
	go logger.run()
}

func (logger *LogHandle) stop() {
	logger.done <- struct{}{}
	logger.flush()
	logger.Close()
}

func (logger *LogHandle) flush() {
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
		logger.files[value] = logger.createFile(value)
	}
	return logger.files[value]
}

func (logger *LogHandle) createFile(value string) *os.File {
	os.MkdirAll(filepath.Join(logger.config.BasePath, value), os.ModeDir)
	file, err := os.Create(filepath.Join(logger.config.BasePath, value, strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%v.csv", time.Now()), " ", "_"), ":", "-")))
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

	if string(payload) == "enable" && !logger.running {
		logger.start()
	} else if string(payload) == "disable" && logger.running {
		logger.stop()
	} else if string(payload) != "enable" && string(payload) != "disable" {
		http.Error(w, "failed to update logger state", http.StatusBadRequest)
		log.Fatalf("log handle: HandleRequest: unexpected body %s\n", payload)
		return
	} else {
		http.Error(w, "failed to update logger state", http.StatusConflict)
		return
	}

	w.Write([]byte{})
}

func (logger *LogHandle) Close() {
	for _, file := range logger.files {
		file.Close()
	}
	logger.files = make(map[string]*os.File, len(logger.files))
}
