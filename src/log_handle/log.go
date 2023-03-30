package log_handle

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/log_handle/models"
	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
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
	trace.Debug().Msg("get log handle")
	return logger
}

func initLogger() {
	trace.Info().Msg("init log handle")
	logger = &LogHandle{
		buffer:      make(map[string][]models.Value),
		autosave:    time.NewTicker(AUTOSAVE_DELAY),
		files:       make(map[string]*os.File),
		done:        make(chan struct{}),
		updates:     make(chan vehicle_models.Update, UPDATE_CHAN_BUF),
		running:     false,
		logSession:  "",
		sendMessage: defaultSendMessage,
		trace:       trace.With().Str("component", LOG_HANDLE_NAME).Logger(),
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
	trace       zerolog.Logger
}

func (logger *LogHandle) UpdateMessage(topic string, payload json.RawMessage, source string) {
	logger.trace.Debug().Str("topic", topic).Str("source", source).Msg("update message")
	if logger.logSession == "" || source == logger.logSession {
		var enable bool
		if err := json.Unmarshal(payload, &enable); err != nil {
			logger.trace.Error().Stack().Err(err).Msg("")
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
		logger.trace.Debug().Str("session", logger.logSession).Msg("set log session")
	} else {
		logger.trace.Warn().Str("source", source).Msg("tried to change running log session")
	}

	logger.notifyState()
}

func (logger *LogHandle) notifyState() {
	logger.trace.Trace().Bool("running", logger.running).Msg("notify state")
	if err := logger.sendMessage("logger/enable", logger.running); err != nil {
		logger.trace.Error().Stack().Err(err).Msg("")
	}
}

func (logger *LogHandle) handleEnable(enable bool) {
	logger.trace.Trace().Bool("enable", enable).Msg("handle enable")
	if enable {
		logger.start()
	} else {
		logger.stop()
	}
}

func (logger *LogHandle) SetSendMessage(sendMessage func(topic string, payload any, target ...string) error) {
	logger.trace.Debug().Msg("set message")
	logger.sendMessage = sendMessage
}

func (logger *LogHandle) HandlerName() string {
	return LOG_HANDLE_NAME
}

func (logger *LogHandle) run() {
	logger.running = true
	defer func() { logger.running = false }()
	logger.trace.Info().Msg("run")
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
			logger.trace.Info().Msg("run stop")
			return
		}
	}
}

func (logger *LogHandle) checkDump() {
	logger.trace.Trace().Msg("check dump")
	for _, buf := range logger.buffer {
		if len(buf) > int(DUMP_SIZE/ROW_SIZE) {
			logger.flush()
			break
		}
	}
}

func (logger *LogHandle) Update(update vehicle_models.Update) {
	if logger.running {
		logger.trace.Trace().Uint16("id", update.ID).Msg("update")
		logger.updates <- update
	}
}

func (logger *LogHandle) start() {
	logger.trace.Debug().Msg("start logger")
	logger.buffer = make(map[string][]models.Value)
	go logger.run()
}

func (logger *LogHandle) stop() {
	logger.trace.Debug().Msg("stop logger")
	logger.done <- struct{}{}
	logger.flush()
	logger.Close()
}

func (logger *LogHandle) flush() {
	logger.trace.Info().Msg("flush")
	for value, buffer := range logger.buffer {
		logger.writeCSV(value, buffer)
	}
	logger.buffer = make(map[string][]models.Value)
}

func (logger *LogHandle) writeCSV(valueName string, buffer []models.Value) {
	file := logger.getFile(valueName)
	data := ""
	for _, value := range buffer {
		logger.trace.Trace().Str("name", valueName).Any("value", value).Msg("write value")
		data += fmt.Sprintf("%d,\"%v\"\n", value.Timestamp.Nanosecond(), value.Value)
	}

	_, err := file.WriteString(data)
	if err != nil {
		logger.trace.Fatal().Stack().Err(err).Msg("")
		return
	}
}

func (logger *LogHandle) getFile(valueName string) *os.File {
	if _, ok := logger.files[valueName]; !ok {
		logger.files[valueName] = logger.createFile(valueName)
	}
	logger.trace.Trace().Str("name", valueName).Msg("get file")
	return logger.files[valueName]
}

func (logger *LogHandle) createFile(valueName string) *os.File {
	err := os.MkdirAll(filepath.Join(LOG_HANDLE_BASE_PATH, valueName), os.ModeDir)
	if err != nil {
		logger.trace.Fatal().Stack().Err(err).Msg("")
		return nil
	}

	path := filepath.Join(LOG_HANDLE_BASE_PATH, valueName, strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%v.csv", time.Now()), " ", "_"), ":", "-"))
	file, err := os.Create(path)
	if err != nil {
		logger.trace.Fatal().Stack().Err(err).Msg("")
		return nil
	}

	logger.trace.Debug().Str("name", valueName).Str("path", path).Msg("create file")
	return file
}

func (logger *LogHandle) Close() error {
	logger.trace.Info().Msg("close")

	var err error
	for _, file := range logger.files {
		if fileErr := file.Close(); err != nil {
			logger.trace.Error().Stack().Err(fileErr).Msg("")
			err = fileErr
		}
	}
	logger.files = make(map[string]*os.File, len(logger.files))
	return err
}

func defaultSendMessage(string, any, ...string) error {
	return errors.New("logger must be registered before using")
}
