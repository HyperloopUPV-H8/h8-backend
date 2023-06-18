package data_transfer

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/common/observable"
	"github.com/HyperloopUPV-H8/Backend-H8/update_factory/models"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const (
	DataTransferHandlerName = "dataTransfer"
	UpdateTopic             = "podData/update"
)

type DataTransfer struct {
	bufMx            *sync.Mutex
	updateBuf        map[uint16]models.Update
	updateObservable observable.ReplayObservable[map[uint16]models.Update]
	ticker           *time.Ticker
	updateTopic      string
	sendMessage      func(topic string, payload any, target ...string) error
	trace            zerolog.Logger
}
type DataTransferTopics struct {
	Update string
}
type DataTransferConfig struct {
	Fps    uint
	Topics DataTransferTopics
}

func New(config DataTransferConfig) DataTransfer {
	trace.Info().Msg("create data transfer")

	dataTransfer := DataTransfer{
		bufMx:            &sync.Mutex{},
		updateBuf:        make(map[uint16]models.Update),
		updateObservable: observable.NewReplayObservable(make(map[uint16]models.Update)),
		ticker:           time.NewTicker(time.Second / time.Duration(config.Fps)),
		updateTopic:      config.Topics.Update,
		sendMessage:      defaultSendMessage,
		trace:            trace.With().Str("component", DataTransferHandlerName).Logger(),
	}

	return dataTransfer
}

func (dataTransfer *DataTransfer) UpdateMessage(topic string, payload json.RawMessage, source string) {
	dataTransfer.trace.Info().Str("source", source).Str("topic", topic).Msg("got message")

	observable.HandleSubscribe[map[uint16]models.Update](&dataTransfer.updateObservable, source, payload,
		func(v map[uint16]models.Update, id string) error {
			err := dataTransfer.sendMessage(UpdateTopic, v, id)

			if websocket.IsCloseError(err) {
				return err
			}

			return nil
		})
}

func (dataTransfer *DataTransfer) SetSendMessage(sendMessage func(topic string, payload any, target ...string) error) {
	dataTransfer.trace.Debug().Msg("set send message")
	dataTransfer.sendMessage = sendMessage
}

func (DataTransfer *DataTransfer) HandlerName() string {
	return DataTransferHandlerName
}

func (dataTransfer *DataTransfer) Run() {
	dataTransfer.trace.Info().Msg("run")
	for {
		<-dataTransfer.ticker.C
		dataTransfer.trySend()
	}
}

func (dataTransfer *DataTransfer) trySend() {
	dataTransfer.bufMx.Lock()
	defer dataTransfer.bufMx.Unlock()

	if len(dataTransfer.updateBuf) == 0 {
		return
	}

	dataTransfer.sendBuf()
}

func (dataTransfer *DataTransfer) sendBuf() {
	dataTransfer.trace.Trace().Msg("send buffer")
	dataTransfer.updateObservable.Next(dataTransfer.updateBuf)
	dataTransfer.updateBuf = make(map[uint16]models.Update, len(dataTransfer.updateBuf))
}

func (dataTransfer *DataTransfer) Update(update models.Update) {
	dataTransfer.bufMx.Lock()
	defer dataTransfer.bufMx.Unlock()
	dataTransfer.trace.Trace().Uint16("id", update.Id).Msg("update")
	dataTransfer.updateBuf[update.Id] = update
}

func defaultSendMessage(string, any, ...string) error {
	return errors.New("data transfer must be registered before use")
}
