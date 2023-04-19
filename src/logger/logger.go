package logger

import (
	"encoding/json"
	"errors"

	"github.com/HyperloopUPV-H8/Backend-H8/logger/data_logger"
	"github.com/HyperloopUPV-H8/Backend-H8/logger/message_logger"
	update_factory_models "github.com/HyperloopUPV-H8/Backend-H8/update_factory/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const (
	LOG_HANDLE_HANDLER_NAME = "logHandle"
)

type Logger struct {
	data    *data_logger.DataLogger
	message *message_logger.MessageLogger

	config Config

	currentSession string

	sendMessage func(topic string, payload any, target ...string) error

	trace zerolog.Logger
}

type EnableMsg struct {
	Client string
	Enable bool
}

func New(config Config) Logger {
	trace.Info().Msg("new log handle")

	messageLogger, err := message_logger.New(config.MessageLogger)
	if err != nil {
		trace.Fatal().Stack().Err(err).Msg("construct message logger")

	}

	return Logger{
		data:    data_logger.New(config.DataLogger),
		message: messageLogger,

		currentSession: "",

		config: config,

		sendMessage: defaultSendMessage,

		trace: trace.With().Str("component", LOG_HANDLE_HANDLER_NAME).Logger(),
	}
}

func (logger *Logger) UpdateMessage(topic string, payload json.RawMessage, source string) {
	logger.trace.Debug().Str("topic", topic).Str("source", source).Msg("update message")
	switch topic {
	case logger.config.Topics.Enable:
		var enable bool
		if err := json.Unmarshal(payload, &enable); err != nil {
			logger.trace.Error().Stack().Err(err).Msg("")
			return
		}

		logger.updateState(enable, source)
	}

	logger.notifyState()
}

func (logger *Logger) updateState(enable bool, source string) {
	logger.trace.Debug().Bool("enable", enable).Str("source", source).Msg("update state")

	if !logger.verifySession(source) {
		logger.trace.Warn().Str("source", source).Msg("tried to change running log session")
		return
	}

	if enable {
		logger.start()
	} else {
		logger.stop()
	}

}

func (logger *Logger) verifySession(session string) bool {
	if logger.currentSession == "" {
		logger.currentSession = session
	}

	return logger.currentSession == session
}

func (logger *Logger) start() {
	logger.data.Start()
}

func (logger *Logger) stop() {
	logger.data.Stop()
	logger.resetSession()
	logger.data.Flush()
}

func (logger *Logger) resetSession() {
	logger.currentSession = ""
}

func (logger *Logger) notifyState() {
	if err := logger.sendMessage(logger.config.Topics.State, logger.data.IsRunning()); err != nil {
		logger.trace.Error().Stack().Err(err).Msg("")
	}
}

func (logger *Logger) NotifyDisconnect(session string) {
	logger.trace.Debug().Str("session", session).Msg("notify disconnect")
	if logger.verifySession(session) {
		logger.stop()
	}
}

func (logger *Logger) UpdateData(data update_factory_models.Update) {
	logger.trace.Trace().Msg("update data")
	logger.data.Update(data)
}

func (logger *Logger) UpdateMsg(msg string) {
	logger.trace.Trace().Msg("update message")
	logger.message.Update(msg)
}

func (logger *Logger) SetSendMessage(sendMessage func(topic string, payload any, target ...string) error) {
	logger.trace.Debug().Msg("set message")
	logger.sendMessage = sendMessage
}

func (logger *Logger) Flush() (err error) {
	logger.trace.Debug().Msg("flush")
	err = logger.data.Flush()
	if err != nil {
		logger.trace.Error().Stack().Err(err).Msg("flush data logger")
	}

	err = logger.message.Flush()
	if err != nil {
		logger.trace.Error().Stack().Err(err).Msg("flush message logger")
	}

	return err
}

func (logger *Logger) Close() (err error) {
	logger.trace.Debug().Msg("close")
	err = logger.data.Close()
	if err != nil {
		logger.trace.Error().Stack().Err(err).Msg("close data logger")
	}

	err = logger.message.Close()
	if err != nil {
		logger.trace.Error().Stack().Err(err).Msg("close message logger")
	}

	return err
}

func (logger *Logger) HandlerName() string {
	return LOG_HANDLE_HANDLER_NAME
}

func defaultSendMessage(string, any, ...string) error {
	return errors.New("logger must be registered before using")
}
