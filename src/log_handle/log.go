package log_handle

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/log_handle/models"
	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

// TODO: Remove hard-coded values
const (
	LOG_HANDLE_NAME      = "logHandle"
	LOG_HANDLE_BASE_PATH = "./log"
	DUMP_SIZE            = 7000
	ROW_SIZE             = 20
	AUTOSAVE_DELAY       = time.Minute
	UPDATE_CHAN_BUF      = 100
)

var (
	logger *LogHandle
)

func Get() *LogHandle {
	if logger == nil {
		initLogger()
	}
	return logger
}

func initLogger() {
	logger = &LogHandle{
		buffer:      make(map[string][]models.Value),
		autosave:    time.NewTicker(AUTOSAVE_DELAY),
		files:       make(map[string]*os.File),
		done:        make(chan struct{}),
		updates:     make(chan vehicle_models.Update, UPDATE_CHAN_BUF),
		running:     false,
		logSession:  "",
		sendMessage: defaultSendMessage,
	}
}

type LogHandle struct {
	buffer      map[string][]models.Value
	autosave    *time.Ticker
	files       map[string]*os.File
	done        chan struct{}
	updates     chan vehicle_models.Update
	running     bool
	logSession  string
	sendMessage func(topic string, payload any, target ...string) error
}

func (logger *LogHandle) UpdateMessage(topic string, payload json.RawMessage, source string) {
	if logger.logSession == "" || source == logger.logSession {
		var enable bool
		if err := json.Unmarshal(payload, &enable); err != nil {
			log.Printf("LogHandle: UpdateMessage: sendMessage: %s\n", err)
			logger.notifyState()
			return
		}

		logger.handleEnable(enable)

		// This can cause locks if the client managing the session disconnects. We should talk how this should work
		if logger.running {
			logger.logSession = source
		} else {
			logger.logSession = ""
		}
	}

	logger.notifyState()
}

func (logger *LogHandle) notifyState() {
	if err := logger.sendMessage("logger/enable", logger.running); err != nil {
		log.Printf("LogHandle: UpdateMessage: sendMessage: %s\n", err)
	}
}

func (logger *LogHandle) handleEnable(enable bool) {
	if enable {
		logger.start()
	} else {
		logger.stop()
	}
}

func (logger *LogHandle) SetSendMessage(sendMessage func(topic string, payload any, target ...string) error) {
	logger.sendMessage = sendMessage
}

func (logger *LogHandle) HandlerName() string {
	return LOG_HANDLE_NAME
}

func (logger *LogHandle) run() {
	logger.running = true
	defer func() { logger.running = false }()
	for {
		select {
		case update := <-logger.updates:
			for name, value := range update.Fields {
				logger.buffer[name] = append(logger.buffer[name], models.Value{
					Value:     value,
					Timestamp: time.Now(),
				})
			}

			logger.checkDump()
		case <-logger.autosave.C:
			logger.flush()
		case <-logger.done:
			return
		}
	}
}

func (logger *LogHandle) checkDump() {
	for _, buf := range logger.buffer {
		if len(buf) > int(DUMP_SIZE/ROW_SIZE) {
			logger.flush()
			break
		}
	}
}

func (logger *LogHandle) Update(update vehicle_models.Update) {
	if logger.running {
		logger.updates <- update
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
	err := os.MkdirAll(filepath.Join(LOG_HANDLE_BASE_PATH, value), os.ModeDir)
	if err != nil {
		log.Fatalf("LogHandle: createFile: %s\n", err)
	}
	path := filepath.Join(LOG_HANDLE_BASE_PATH, value, strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%v.csv", time.Now()), " ", "_"), ":", "-"))
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("LogHandle: WriteCSV: %s\n", err)
	}
	return file
}

func (logger *LogHandle) Close() {
	for _, file := range logger.files {
		file.Close()
	}
	logger.files = make(map[string]*os.File, len(logger.files))
}

func defaultSendMessage(string, any, ...string) error {
	return errors.New("logger must be registered before using")
}
