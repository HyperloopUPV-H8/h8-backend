package old_logger

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const LOG_HANDLE_HANDLER_NAME = "logHandle"

type Logger struct {
	subloggers map[string]SubLogger

	currentSession string
	isRunning      bool

	config Config

	sendMessage func(topic string, payload any, target ...string) error

	trace zerolog.Logger
}

func New(subloggers map[string]SubLogger, config Config) Logger {
	trace.Info().Msg("new log handle")

	return Logger{
		subloggers: subloggers,

		currentSession: "",
		isRunning:      false,

		config: config,

		sendMessage: defaultSendMessage,

		trace: trace.With().Str("component", LOG_HANDLE_HANDLER_NAME).Logger(),
	}
}

func (logger *Logger) UpdateMessage(topic string, payload json.RawMessage, source string) {
	logger.trace.Debug().Str("topic", topic).Str("source", source).Msg("update message")
	switch topic {
	case logger.config.Topics.Enable:
		logger.handleEnable(payload, source)
	}

	logger.notifyState()
}

func (logger *Logger) handleEnable(payload json.RawMessage, source string) error {
	var enable bool
	err := json.Unmarshal(payload, &enable)
	if err != nil {
		logger.trace.Error().Stack().Err(err).Msg("unmarshal enable")
		return err
	}

	err = logger.updateState(enable, source)
	return err
}

func (logger *Logger) updateState(enable bool, source string) error {
	logger.trace.Debug().Bool("enable", enable).Str("source", source).Msg("update state")

	if !logger.verifySession(source) {
		logger.trace.Warn().Str("source", source).Msg("tried to change running log session")
		return fmt.Errorf("%s tried to change running log session of %s", source, logger.currentSession)
	}

	var err error
	if enable {
		err = logger.start()
	} else {
		err = logger.stop()
	}

	return err
}

func (logger *Logger) verifySession(session string) bool {
	if logger.currentSession == "" {
		logger.currentSession = session
	}

	return logger.currentSession == session
}

func (logger *Logger) start() (err error) {
	for _, sublogger := range logger.subloggers {
		startErr := sublogger.Start()
		if startErr != nil {
			err = startErr
		}
	}

	logger.isRunning = true
	return err
}

func (logger *Logger) stop() (err error) {
	for _, sublogger := range logger.subloggers {
		stopErr := sublogger.Stop()
		if stopErr != nil {
			err = stopErr
		}
	}

	logger.isRunning = false
	logger.resetSession()
	flushErr := logger.Flush()
	if flushErr != nil {
		err = flushErr
	}
	return err
}

func (logger *Logger) resetSession() {
	logger.currentSession = ""
}

func (logger *Logger) notifyState() {
	if err := logger.sendMessage(logger.config.Topics.State, logger.isRunning); err != nil {
		logger.trace.Error().Stack().Err(err).Msg("")
	}
}

func (logger *Logger) NotifyDisconnect(session string) {
	logger.trace.Debug().Str("session", session).Msg("notify disconnect")
	if logger.verifySession(session) {
		logger.stop()
	}
}

func (logger *Logger) SetSendMessage(sendMessage func(topic string, payload any, target ...string) error) {
	logger.trace.Debug().Msg("set message")
	logger.sendMessage = sendMessage
}

func (logger *Logger) Flush() (err error) {
	logger.trace.Debug().Msg("flush")

	for _, sublogger := range logger.subloggers {
		flushErr := sublogger.Flush()
		if flushErr != nil {
			err = flushErr
		}
	}

	return err
}

func (logger *Logger) Close() (err error) {
	for _, sublogger := range logger.subloggers {
		closeErr := sublogger.Close()
		if closeErr != nil {
			err = closeErr
		}
	}

	return err
}

func (logger *Logger) HandlerName() string {
	return LOG_HANDLE_HANDLER_NAME
}

func defaultSendMessage(string, any, ...string) error {
	return errors.New("logger must be registered before using")
}
