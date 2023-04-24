package logger_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	sendMessage    func(topic string, payload any, target ...string) error
	config         Config
	trace          zerolog.Logger
}

func NewLoggerHandler(loggers map[string]Logger, config Config) LoggerHandler {
	trace.Info().Msg("new LoggerHandler")

	os.MkdirAll(config.BasePath, os.ModeDir)

	return LoggerHandler{
		loggers:        loggers,
		loggableChan:   make(chan Loggable),
		currentSession: "",
		isRunning:      false,
		sendMessage:    defaultSendMessage,
		config:         config,
		trace:          trace.With().Str("component", LOG_HANDLER_HANDLER_NAME).Logger(),
	}
}

func (handler *LoggerHandler) Log(loggable Loggable) {
	if handler.isRunning {
		handler.loggableChan <- loggable
	}
}

func (handler *LoggerHandler) UpdateMessage(topic string, payload json.RawMessage, source string) {
	handler.trace.Debug().Str("topic", topic).Str("source", source).Msg("update message")
	switch topic {
	case handler.config.Topics.Enable:
		handler.handleEnable(payload, source)
	}

	handler.notifyState()
}

func (handler *LoggerHandler) handleEnable(payload json.RawMessage, source string) error {
	var enable bool
	err := json.Unmarshal(payload, &enable)
	if err != nil {
		handler.trace.Error().Stack().Err(err).Msg("unmarshal enable")
		return err
	}

	handler.updateState(enable, source)
	return nil
}

func (handler *LoggerHandler) updateState(enable bool, source string) error {
	handler.trace.Debug().Bool("enable", enable).Str("source", source).Msg("update state")

	if !handler.verifySession(source) {
		handler.trace.Warn().Str("source", source).Msg("tried to change running log session")
		return fmt.Errorf("%s tried to change running log session of %s", source, handler.currentSession)
	}

	if enable && !handler.isRunning {
		handler.start()
	} else if !enable && handler.isRunning {
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
	sessionDirName := strconv.Itoa(int(time.Now().Unix()))
	path := filepath.Join(handler.config.BasePath, sessionDirName)
	os.MkdirAll(path, os.ModeDir)

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

	for _, logger := range activeLoggers {
		close(logger.Input) // loggers should clean up after channel has been closed
	}
}

func (handler *LoggerHandler) stop() {
	handler.trace.Info().Str("logger session", handler.currentSession).Msg("Stoped logging")
	close(handler.loggableChan) // triggers loggers clean-up
	handler.isRunning = false
	handler.currentSession = ""
}

func (handler *LoggerHandler) NotifyDisconnect(session string) {
	handler.trace.Debug().Str("session", session).Msg("notify disconnect")
	if handler.verifySession(session) && handler.isRunning {
		handler.stop()
	}
}

func (handler *LoggerHandler) notifyState() {
	if err := handler.sendMessage(handler.config.Topics.State, handler.isRunning); err != nil {
		handler.trace.Error().Stack().Err(err).Msg("")
	}
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
