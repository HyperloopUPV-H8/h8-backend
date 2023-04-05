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

const (
	LOG_HANDLE_HANDLER_NAME = "logHandle"
	UPDATE_CHAN_BUF         = 100
)

var (
	loggerConfig = LoggerConfig{
		DumpSize:      7000,
		RowSize:       20,
		AutosaveDelay: "1m",
		Path:          "./log",
		Topics: LoggerTopics{
			Enable: "logger/enable",
			State:  "logger/state",
		},
	}
	logger *LogHandle
)

type LoggerConfig struct {
	DumpSize      uint   `toml:"dump_size"`
	RowSize       uint   `toml:"row_size"`
	AutosaveDelay string `toml:"autosave_delay"`
	Path          string
	Topics        LoggerTopics
}

type LoggerTopics struct {
	Enable string
	State  string
}

func SetConfig(config LoggerConfig) {
	loggerConfig = config
}

func Get() *LogHandle {
	if logger == nil {
		initLogger()
	}
	trace.Debug().Msg("get log handle")
	return logger
}

func initLogger() {
	trace.Info().Msg("init log handle")

	autosaveDelay, err := time.ParseDuration(loggerConfig.AutosaveDelay)
	if err != nil {
		trace.Fatal().Stack().Err(err).Str("LOGGER_AUTOSAVE_DELAY", loggerConfig.AutosaveDelay).Msg("")
	}

	logger = &LogHandle{
		files:    make(map[string]*os.File),
		buffer:   make(map[string][]models.Value),
		autosave: time.NewTicker(autosaveDelay),

		done:    make(chan struct{}),
		updates: make(chan vehicle_models.Update, UPDATE_CHAN_BUF),

		isRunning: false,
		session:   "",

		topics:      LoggerTopics{Enable: loggerConfig.Topics.Enable, State: loggerConfig.Topics.State},
		sendMessage: defaultSendMessage,

		path:          loggerConfig.Path,
		dumpThreshold: loggerConfig.DumpSize / loggerConfig.RowSize,

		trace: trace.With().Str("component", LOG_HANDLE_HANDLER_NAME).Logger(),
	}
}

type LogHandle struct {
	files    map[string]*os.File
	buffer   map[string][]models.Value
	autosave *time.Ticker

	done    chan struct{}
	updates chan vehicle_models.Update

	isRunning bool
	session   string

	topics struct {
		Enable string
		State  string
	}
	sendMessage func(topic string, payload any, target ...string) error

	path          string
	dumpThreshold uint

	trace zerolog.Logger
}

func (logger *LogHandle) UpdateMessage(topic string, payload json.RawMessage, source string) {
	logger.trace.Debug().Str("topic", topic).Str("source", source).Msg("update message")
	switch topic {
	case logger.topics.Enable:
		logger.handleEnableRequest(topic, payload, source)
	}
	logger.notifyState()
}

func (logger *LogHandle) handleEnableRequest(topic string, payload json.RawMessage, source string) {
	if logger.session != "" && source != logger.session {
		logger.trace.Warn().Str("source", source).Msg("tried to change running log session")
		return
	}

	var enable bool
	if err := json.Unmarshal(payload, &enable); err != nil {
		logger.trace.Error().Stack().Err(err).Msg("")
		return
	}

	logger.handleEnable(enable)

	// This can cause locks if the client managing the session disconnects. We should talk how this should work
	if logger.isRunning {
		logger.session = source
	} else {
		logger.session = ""
	}

	logger.trace.Debug().Str("session", logger.session).Msg("set log session")
}

func (logger *LogHandle) notifyState() {
	logger.trace.Trace().Bool("running", logger.isRunning).Msg("notify state")
	if err := logger.sendMessage(logger.topics.State, logger.isRunning); err != nil {
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
	return LOG_HANDLE_HANDLER_NAME
}

func (logger *LogHandle) run() {
	logger.isRunning = true
	defer func() { logger.isRunning = false }()
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
		if len(buf) > int(logger.dumpThreshold) {
			logger.flush()
			break
		}
	}
}

func (logger *LogHandle) Update(update vehicle_models.Update) {
	if logger.isRunning {
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
	basePath := logger.path
	err := os.MkdirAll(filepath.Join(basePath, valueName), os.ModeDir)
	if err != nil {
		logger.trace.Fatal().Stack().Err(err).Msg("")
		return nil
	}

	path := filepath.Join(basePath, valueName, strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%v.csv", time.Now()), " ", "_"), ":", "-"))
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
