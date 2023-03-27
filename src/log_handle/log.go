package log_handle

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/log_handle/models"
	ws_models "github.com/HyperloopUPV-H8/Backend-H8/websocket_handle/models"
)

type LogHandle struct {
	flushing   bool
	buffer     map[string][]models.Value
	config     models.Config
	fileMx     sync.Mutex
	files      map[string]*os.File
	done       chan struct{}
	running    bool
	logSession string
	channel    chan ws_models.MessageTarget
}

func NewLogger(config models.Config) (*LogHandle, chan ws_models.MessageTarget) {
	logHandle := &LogHandle{
		flushing: false,
		buffer:   make(map[string][]models.Value),
		config:   config,
		fileMx:   sync.Mutex{},
		files:    make(map[string]*os.File),
		done:     make(chan struct{}),
		channel:  make(chan ws_models.MessageTarget),
	}

	go logHandle.listenWS()

	return logHandle, logHandle.channel
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
	log.Println("Starting logger")
	logger.buffer = make(map[string][]models.Value)
	go logger.run()
}

func (logger *LogHandle) stop() {
	log.Println("Stopping logger")
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
	_, err := file.WriteString(data)
	if err != nil {
		log.Fatalf("LogHandle: writeCSV: %s\n", err)
	}
}

func (logger *LogHandle) getFile(value string) *os.File {
	if _, ok := logger.files[value]; !ok {
		logger.files[value] = logger.createFile(value)
	}
	return logger.files[value]
}

func (logger *LogHandle) createFile(value string) *os.File {
	err := os.MkdirAll(filepath.Join(logger.config.BasePath, value), os.ModeDir)
	if err != nil {
		log.Fatalf("LogHandle: createFile: %s\n", err)
	}
	path := filepath.Join(logger.config.BasePath, value, strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%v.csv", time.Now()), " ", "_"), ":", "-"))
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("LogHandle: WriteCSV: %s\n", err)
	}
	return file
}

func (logger *LogHandle) listenWS() {
	for msg := range logger.channel {
		log.Println(msg)
		if logger.logSession == "" || msg.Target[0] == logger.logSession {
			var enable bool
			err := json.Unmarshal(msg.Msg.Msg, &enable)
			if err != nil {
				log.Printf("logger: listenWS: %s\n", err)
				continue
			}
			if enable {
				logger.start()
				logger.logSession = msg.Target[0]
			} else {
				logger.stop()
				logger.logSession = ""
			}

			logger.channel <- ws_models.NewMessageTargetRaw([]string{}, "logger", msg.Msg.Msg)
		}
	}
}

func (logger *LogHandle) Close() {
	for _, file := range logger.files {
		file.Close()
	}
	logger.files = make(map[string]*os.File, len(logger.files))
}
