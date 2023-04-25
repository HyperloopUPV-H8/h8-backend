package data_transfer

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/update_factory/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const (
	DATA_TRANSFER_HANDLER_NAME = "dataTransfer"
)

type DataTransfer struct {
	bufMx       *sync.Mutex
	packetBuf   map[uint16]models.Update
	ticker      *time.Ticker
	updateTopic string
	sendMessage func(topic string, payload any, target ...string) error
	trace       zerolog.Logger
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
		bufMx:       &sync.Mutex{},
		packetBuf:   make(map[uint16]models.Update),
		ticker:      time.NewTicker(time.Second / time.Duration(config.Fps)),
		updateTopic: config.Topics.Update,
		sendMessage: defaultSendMessage,
		trace:       trace.With().Str("component", DATA_TRANSFER_HANDLER_NAME).Logger(),
	}

	return dataTransfer
}

func (dataTransfer *DataTransfer) UpdateMessage(topic string, payload json.RawMessage, source string) {
	dataTransfer.trace.Warn().Str("source", source).Str("topic", topic).Msg("got message")
}

func (dataTransfer *DataTransfer) SetSendMessage(sendMessage func(topic string, payload any, target ...string) error) {
	dataTransfer.trace.Debug().Msg("set send message")
	dataTransfer.sendMessage = sendMessage
}

func (DataTransfer *DataTransfer) HandlerName() string {
	return DATA_TRANSFER_HANDLER_NAME
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

	if len(dataTransfer.packetBuf) == 0 {
		return
	}

	dataTransfer.sendBuf()
}

func (dataTransfer *DataTransfer) sendBuf() {
	dataTransfer.trace.Trace().Msg("send buffer")
	if err := dataTransfer.sendMessage(dataTransfer.updateTopic, dataTransfer.packetBuf); err != nil {
		dataTransfer.trace.Error().Stack().Err(err).Msg("")
		return
	}

	dataTransfer.packetBuf = make(map[uint16]models.Update, len(dataTransfer.packetBuf))
}

func (dataTransfer *DataTransfer) Update(update models.Update) {
	dataTransfer.bufMx.Lock()
	defer dataTransfer.bufMx.Unlock()

	dataTransfer.trace.Trace().Uint16("id", update.ID).Msg("update")
	dataTransfer.packetBuf[update.ID] = update
}

func defaultSendMessage(string, any, ...string) error {
	return errors.New("data transfer must be registered before use")
}
