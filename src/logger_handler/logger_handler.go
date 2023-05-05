package logger_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const LOG_HANDLER_HANDLER_NAME = "LoggerHandler"

type LoggerHandler struct {
	loggers        map[string]Logger
	loggableChan   chan Loggable
	currentSession string
	isRunning      bool
	isRunningMx    *sync.Mutex
	sendMessage    func(topic string, payload any, target ...string) error
	config         Config
	trace          zerolog.Logger
}

func NewLoggerHandler(loggers map[string]Logger, config Config) LoggerHandler {
	trace.Info().Msg("new LoggerHandler")

	os.MkdirAll(config.BasePath, os.ModePerm)

	return LoggerHandler{
		loggers:        loggers,
		loggableChan:   make(chan Loggable),
		currentSession: "",
		isRunning:      false,
		isRunningMx:    &sync.Mutex{},

		sendMessage: defaultSendMessage,
		config:      config,
		trace:       trace.With().Str("component", LOG_HANDLER_HANDLER_NAME).Logger(),
	}
}

func (handler *LoggerHandler) Log(loggable Loggable) {
	handler.isRunningMx.Lock()
	if handler.isRunning {
		handler.loggableChan <- loggable
	}
	handler.isRunningMx.Unlock()
}

func (handler *LoggerHandler) UpdateMessage(topic string, payload json.RawMessage, source string) {
	handler.trace.Info().Str("topic", topic).Str("source", source).Msg("update message")
	switch topic {
	case handler.config.Topics.Enable:
		var enable bool
		err := json.Unmarshal(payload, &enable)
		if err != nil {
			handler.trace.Error().Stack().Err(err).Msg("unmarshal enable")
		}

		handler.handleEnable(enable, source)
	}

	handler.notifyState()
}

func (handler *LoggerHandler) handleEnable(enable bool, source string) error {
	if !handler.verifySession(source) {
		handler.trace.Warn().Str("source", source).Msg("tried to change running log session")
		return fmt.Errorf("%s tried to change running log session of %s", source, handler.currentSession)
	}

	handler.isRunningMx.Lock()
	if enable && !handler.isRunning {
		handler.isRunningMx.Unlock()
		handler.start()
	} else if !enable && handler.isRunning {
		handler.isRunningMx.Unlock()
		handler.stop()
	}

	return nil
}

func (handler *LoggerHandler) verifySession(session string) bool {
	if handler.currentSession == "" {
		handler.currentSession = session
	}

	return handler.currentSession == session
}

func (handler *LoggerHandler) start() {
	handler.trace.Info().Str("logger session", handler.currentSession).Msg("Started logging")
	handler.loggableChan = make(chan Loggable)
	currentTime := time.Now()
	sessionDirName := fmt.Sprintf("%d_%d_%d - %d_%dh", currentTime.Day(), currentTime.Month(), currentTime.Year(), currentTime.Hour(), currentTime.Minute())
	path := filepath.Join(handler.config.BasePath, sessionDirName)
	os.MkdirAll(path, os.ModePerm)

	activeLoggers := handler.createActiveLoggers(path)

	go startBroadcastRoutine(activeLoggers, handler.loggableChan)
	handler.isRunning = true
}

func (handler *LoggerHandler) createActiveLoggers(path string) []ActiveLogger {
	activeLoggers := make([]ActiveLogger, 0)

	for _, logger := range handler.loggers {
		inputChan := logger.Start(path)

		activeLoggers = append(activeLoggers, ActiveLogger{
			Ids:   logger.Ids(),
			Input: inputChan,
		})
	}

	return activeLoggers
}

func startBroadcastRoutine(activeLoggers []ActiveLogger, generalInput <-chan Loggable) {
	for loggable := range generalInput {
		for _, logger := range activeLoggers {
			if logger.Ids.Has(loggable.Id()) {
				logger.Input <- loggable
			}
		}
	}

	defer func() {
		for _, logger := range activeLoggers {
			close(logger.Input)
		}
	}()
}

func (handler *LoggerHandler) stop() {
	handler.trace.Info().Str("logger session", handler.currentSession).Msg("Stoped logging")
	handler.isRunning = false
	close(handler.loggableChan) // triggers loggers clean-up
	handler.currentSession = ""
}

func (handler *LoggerHandler) NotifyDisconnect(session string) {
	handler.trace.Debug().Str("session", session).Msg("notify disconnect")
	if handler.verifySession(session) {
		handler.isRunningMx.Lock()
		if handler.isRunning {
			handler.stop()
		}
		handler.isRunningMx.Unlock()

		handler.currentSession = ""
	}
}

func (handler *LoggerHandler) notifyState() error {
	handler.isRunningMx.Lock()
	runningCopy := handler.isRunning
	handler.isRunningMx.Unlock()

	if err := handler.sendMessage(handler.config.Topics.Enable, runningCopy); err != nil {
		handler.trace.Error().Stack().Err(err).Msg("")
		return err
	}

	return nil
}

func (handler *LoggerHandler) SetSendMessage(sendMessage func(topic string, payload any, target ...string) error) {
	handler.trace.Debug().Msg("set message")
	handler.sendMessage = sendMessage
}

func (handler *LoggerHandler) HandlerName() string {
	return LOG_HANDLER_HANDLER_NAME
}

func defaultSendMessage(string, any, ...string) error {
	return errors.New("logger must be registered before using")
}
