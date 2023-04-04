package data_transfer

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const (
	DATA_TRANSFER_HANDLER_NAME = "dataTransfer"
)

var (
	dataTransfer *DataTransfer
)

func Get() *DataTransfer {
	if dataTransfer == nil {
		initDataTransfer()
	}
	trace.Debug().Msg("get data transfer")
	return dataTransfer
}

func initDataTransfer() {
	trace.Info().Msg("init data transfer")

	refreshRate, err := strconv.ParseInt(os.Getenv("DATA_TRANSFER_FPS"), 10, 32)
	if err != nil {
		trace.Fatal().Stack().Err(err).Str("DATA_TRANSFER_FPS", os.Getenv("DATA_TRANSFER_FPS")).Msg("")
	}

	dataTransfer = &DataTransfer{
		bufMx:       &sync.Mutex{},
		packetBuf:   make(map[uint16]models.Update),
		ticker:      time.NewTicker(time.Second / time.Duration(refreshRate)),
		sendMessage: defaultSendMessage,
		trace:       trace.With().Str("component", DATA_TRANSFER_HANDLER_NAME).Logger(),
	}

	go dataTransfer.run()
}

type DataTransfer struct {
	bufMx       *sync.Mutex
	packetBuf   map[uint16]models.Update
	ticker      *time.Ticker
	sendMessage func(topic string, payload any, target ...string) error
	trace       zerolog.Logger
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

func (dataTransfer *DataTransfer) run() {
	dataTransfer.trace.Info().Msg("run")
	for {
		<-dataTransfer.ticker.C
		if len(dataTransfer.packetBuf) == 0 {
			continue
		}

		dataTransfer.sendBuf()
	}
}

func (dataTransfer *DataTransfer) sendBuf() {
	dataTransfer.bufMx.Lock()
	defer dataTransfer.bufMx.Unlock()

	dataTransfer.trace.Trace().Msg("send buffer")
	if err := dataTransfer.sendMessage(os.Getenv("DATA_TRANSFER_TOPIC"), dataTransfer.packetBuf); err != nil {
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
