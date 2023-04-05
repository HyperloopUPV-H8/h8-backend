package connection_transfer

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/HyperloopUPV-H8/Backend-H8/connection_transfer/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const (
	CONNECTION_TRANSFER_HANDLER_NAME = "connectionTransfer"
)

type ConnectionTransferConfig struct {
	UpdateTopic string `toml:"update_topic"`
}

var (
	connectionTransferConfig = ConnectionTransferConfig{
		UpdateTopic: "connection/update",
	}
	connectionTransfer *ConnectionTransfer
)

func SetConfig(config ConnectionTransferConfig) {
	connectionTransferConfig = config
}

func Get() *ConnectionTransfer {
	if connectionTransfer == nil {
		initConnectionTransfer()
	}
	trace.Debug().Msg("get connection transfer")
	return connectionTransfer
}

func initConnectionTransfer() {
	trace.Info().Msg("init connection transfer")
	connectionTransfer = &ConnectionTransfer{
		writeMx:     &sync.Mutex{},
		boardStatus: make(map[string]models.Connection),
		sendMessage: defaultSendMessage,
		updateTopic: connectionTransferConfig.UpdateTopic,
		trace:       trace.With().Str("component", CONNECTION_TRANSFER_HANDLER_NAME).Logger(),
	}
}

type ConnectionTransfer struct {
	writeMx     *sync.Mutex
	boardStatus map[string]models.Connection
	sendMessage func(topic string, payload any, target ...string) error
	updateTopic string
	trace       zerolog.Logger
}

func (connectionTransfer *ConnectionTransfer) UpdateMessage(topic string, payload json.RawMessage, source string) {
	connectionTransfer.trace.Trace().Str("source", source).Str("topic", topic).Msg("got message")
	connectionTransfer.send()
}

func (connectionTransfer *ConnectionTransfer) SetSendMessage(sendMessage func(topic string, payload any, target ...string) error) {
	connectionTransfer.trace.Debug().Msg("set send message")
	connectionTransfer.sendMessage = sendMessage
}

func (connectionTransfer *ConnectionTransfer) HandlerName() string {
	return CONNECTION_TRANSFER_HANDLER_NAME
}

func (connectionTransfer *ConnectionTransfer) Update(name string, up bool) {
	connectionTransfer.writeMx.Lock()
	defer connectionTransfer.writeMx.Unlock()

	connectionTransfer.trace.Debug().Str("connection", name).Bool("isConnected", up).Msg("update connection state")

	connectionTransfer.boardStatus[name] = models.Connection{
		Name:        name,
		IsConnected: up,
	}

	connectionTransfer.send()
}

func (connectionTransfer *ConnectionTransfer) send() {
	connectionTransfer.trace.Debug().Msg("send connections")
	if err := connectionTransfer.sendMessage(connectionTransfer.updateTopic, connectionTransfer.boardStatus); err != nil {
		connectionTransfer.trace.Error().Stack().Err(err).Msg("")
		return
	}
}

func defaultSendMessage(string, any, ...string) error {
	return errors.New("connection transfer must be registered before use")
}
