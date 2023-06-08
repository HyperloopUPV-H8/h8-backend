package message_transfer

import (
	"encoding/json"
	"errors"

	"github.com/HyperloopUPV-H8/Backend-H8/common/observable"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const (
	MessageTransferHandlerName = "messageTransfer"
	UpdateTopic                = "message/update"
)

type MessageTransfer struct {
	updateTopic       string
	messageObservable observable.NoReplayObservable[any]
	sendMessage       func(topic string, payload any, target ...string) error
	trace             zerolog.Logger
}
type MessageTransferConfig struct {
	UpdateTopic string `toml:"update_topic"`
}

func New(config MessageTransferConfig) MessageTransfer {
	trace.Info().Msg("new message transfer")
	return MessageTransfer{
		updateTopic:       config.UpdateTopic,
		messageObservable: observable.NewNoReplayObservable[any](),
		sendMessage:       defaultSendMessage,
		// messageObservable: observable.NewWsObservable[any](struct{}{}, func(v any, id string) error { return defaultSendMessage(id, v) }), //FIXME: change struct{}{}
		trace: trace.With().Str("component", MessageTransferHandlerName).Logger(),
	}
}

func (messageTransfer *MessageTransfer) SendMessage(message any) error {
	// messageTransfer.trace.Warn().Uint16("id", message.ID).Str("type", message.Type).Str("desc", message.Description).Msg("send message")
	messageTransfer.messageObservable.Next(message)
	return nil
}

func (messageTransfer *MessageTransfer) UpdateMessage(topic string, payload json.RawMessage, source string) {
	messageTransfer.trace.Info().Str("source", source).Str("topic", topic).Msg("got message")

	observable.HandleSubscribe[any](&messageTransfer.messageObservable, source, payload,
		func(v any, id string) error {
			return messageTransfer.sendMessage(UpdateTopic, v, id)
		})
}

func (messageTransfer *MessageTransfer) SetSendMessage(sendMessage func(topic string, payload any, targets ...string) error) {
	messageTransfer.trace.Debug().Msg("set send message")
	messageTransfer.sendMessage = sendMessage
}

func (messageTransfer *MessageTransfer) HandlerName() string {
	return MessageTransferHandlerName
}

func defaultSendMessage(string, any, ...string) error {
	return errors.New("message transfer must be registered before using")
}
