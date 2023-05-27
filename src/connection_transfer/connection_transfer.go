package connection_transfer

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/common/observable"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const (
	ConnectionTransferHandlerName = "connectionTransfer"
	UpdateTopic                   = "connection/update"
)

type ConnectionTransfer struct {
	writeMx               *sync.Mutex
	boardStatus           map[string]Connection
	boardStatusObservable observable.ReplayObservable[[]Connection]
	sendMessage           func(topic string, payload any, target ...string) error
	updateTopic           string
	trace                 zerolog.Logger
}

type ConnectionTransferConfig struct {
	UpdateTopic string `toml:"update_topic"`
}

func New(config ConnectionTransferConfig) ConnectionTransfer {
	trace.Info().Msg("new connection transfer")

	return ConnectionTransfer{
		writeMx:               &sync.Mutex{},
		boardStatus:           make(map[string]Connection),
		boardStatusObservable: observable.NewReplayObservable(make([]Connection, 0)),
		sendMessage:           defaultSendMessage,
		updateTopic:           config.UpdateTopic,
		trace:                 trace.With().Str("component", ConnectionTransferHandlerName).Logger(),
	}
}

func (connectionTransfer *ConnectionTransfer) UpdateMessage(topic string, payload json.RawMessage, source string) {
	connectionTransfer.trace.Trace().Str("source", source).Str("topic", topic).Msg("got message")

	var subscribe ConnectionSubscription
	err := json.Unmarshal(payload, &subscribe)

	if err != nil {
		connectionTransfer.trace.Error().Err(err).Msg("unmarshaling payload")
	}

	if subscribe {
		connectionTransfer.AddObserver(source)
	} else {
		connectionTransfer.boardStatusObservable.Unsubscribe(source)
	}

}

func (connectionTransfer *ConnectionTransfer) AddObserver(id string) {
	observer := observable.NewWsObserver(id, func(v []Connection) {
		err := connectionTransfer.sendMessage(UpdateTopic, v, id)

		if err != nil {
			connectionTransfer.boardStatusObservable.Unsubscribe(id)
		}
	})

	connectionTransfer.boardStatusObservable.Subscribe(observer)
}

func (connectionTransfer *ConnectionTransfer) SetSendMessage(sendMessage func(topic string, payload any, target ...string) error) {
	connectionTransfer.trace.Debug().Msg("set send message")
	connectionTransfer.sendMessage = sendMessage
}

func (connectionTransfer *ConnectionTransfer) HandlerName() string {
	return ConnectionTransferHandlerName
}

func (connectionTransfer *ConnectionTransfer) Update(name string, IsConnected bool) {
	connectionTransfer.writeMx.Lock()
	defer connectionTransfer.writeMx.Unlock()

	connectionTransfer.trace.Debug().Str("connection", name).Bool("isConnected", IsConnected).Msg("update connection state")

	connectionTransfer.boardStatus[name] = Connection{
		Name:        name,
		IsConnected: IsConnected,
	}

	connectionTransfer.boardStatusObservable.Next(common.Values(connectionTransfer.boardStatus))
}

func defaultSendMessage(string, any, ...string) error {
	return errors.New("connection transfer must be registered before use")
}
