package logger_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const LOG_HANDLER_HANDLER_NAME = "LoggerHandler"

type LoggerHandler struct {
	loggers        map[string]Logger
	loggableChan   <-chan Loggable
	currentSession string
	isRunning      bool
	sendMessage    func(topic string, payload any, target ...string) error
	config         Config
	trace          zerolog.Logger
}

func NewLoggerHandler(loggers map[string]Logger, loggableChan <-chan Loggable, config Config) LoggerHandler {
	trace.Info().Msg("new LoggerHandler")

	os.Mkdir("log", os.ModeDir) //TODO: comprobar si borra la carpeta

	return LoggerHandler{
		loggers:        loggers,
		loggableChan:   loggableChan,
		currentSession: "",
		isRunning:      false,
		sendMessage:    defaultSendMessage,
		config:         config,
		trace:          trace.With().Str("component", LOG_HANDLER_HANDLER_NAME).Logger(),
	}
}

func (handler *LoggerHandler) Listen() {
	for loggable := range handler.loggableChan {
		//TODO: send only to interested logger instead of all
		for _, logger := range handler.loggers {
			logger.Log(loggable)
		}
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

	err = handler.updateState(enable, source)
	return err
}

func (handler *LoggerHandler) updateState(enable bool, source string) error {
	handler.trace.Debug().Bool("enable", enable).Str("source", source).Msg("update state")

	if !handler.verifySession(source) {
		handler.trace.Warn().Str("source", source).Msg("tried to change running log session")
		return fmt.Errorf("%s tried to change running log session of %s", source, handler.currentSession)
	}

	var err error
	if enable {
		err = handler.start()
	} else {
		err = handler.stop()
	}

	return err
}

func (handler *LoggerHandler) verifySession(session string) bool {
	if handler.currentSession == "" {
		handler.currentSession = session
	}

	return handler.currentSession == session
}

func (handler *LoggerHandler) start() (err error) {
	for _, logger := range handler.loggers {
		startErr := logger.Start()
		if startErr != nil {
			err = startErr
		}
	}

	handler.isRunning = true
	return err
}

func (handler *LoggerHandler) stop() (err error) {
	for _, logger := range handler.loggers {
		stopErr := logger.Stop()
		if stopErr != nil {
			err = stopErr
		}
	}

	handler.isRunning = false
	handler.currentSession = ""
	flushErr := handler.Flush()
	if flushErr != nil {
		err = flushErr
	}
	return err
}

func (handler *LoggerHandler) Log(loggable Loggable) {
	for _, logger := range handler.loggers {
		logger.Log(loggable)
	}
}

func (handler *LoggerHandler) Flush() (err error) {
	handler.trace.Debug().Msg("flush")

	for _, logger := range handler.loggers {
		flushErr := logger.Flush()
		if flushErr != nil {
			err = flushErr
		}
	}

	return err
}

func (handler *LoggerHandler) Close() (err error) {
	for _, logger := range handler.loggers {
		closeErr := logger.Close()
		if closeErr != nil {
			err = closeErr
		}
	}

	return err
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
