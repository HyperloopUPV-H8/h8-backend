package blcu

import (
	"encoding/json"
	"net"

	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

type BLCU struct {
	addr       net.TCPAddr
	boardToId  map[string]uint16
	ackChannel chan struct{}

	sendOrder   func(models.Order) error
	sendMessage func(topic string, payload any, targets ...string) error

	config BLCUConfig

	trace zerolog.Logger
}

func NewBLCU(laddr net.TCPAddr, boardIds map[string]uint16, config BLCUConfig) BLCU {
	trace.Info().Msg("New BLCU")

	return BLCU{
		addr:        laddr,
		boardToId:   boardIds,
		ackChannel:  make(chan struct{}, BLCU_ACK_CHAN_BUF),
		trace:       trace.With().Str("component", BLCU_COMPONENT_NAME).Logger(),
		config:      config,
		sendOrder:   func(o models.Order) error { return nil },
		sendMessage: func(topic string, payload any, targets ...string) error { return nil },
	}
}

func (blcu *BLCU) HandlerName() string {
	return BLCU_HANDLER_NAME
}

func (blcu *BLCU) SetSendMessage(sendMessage func(topic string, payload any, targets ...string) error) {
	blcu.trace.Debug().Msg("Set send message")
	blcu.sendMessage = sendMessage
}

func (blcu *BLCU) SetSendOrder(sendOrder func(o models.Order) error) {
	blcu.sendOrder = sendOrder
}

func (blcu *BLCU) UpdateMessage(topic string, payload json.RawMessage, source string) {
	blcu.trace.Debug().Str("topic", topic).Str("source", source).Msg("Update message")
	switch topic {
	case blcu.config.Topics.Upload:
		err := blcu.upload(payload)
		if err != nil {
			blcu.notifyUploadFailure()
		} else {
			blcu.notifyUploadSuccess()
		}

	case blcu.config.Topics.Download:
		board, data, err := blcu.download(payload)
		if err != nil {
			blcu.notifyDownloadFailure()
		} else {
			blcu.notifyDownloadSuccess(data)
			blcu.writeDownloadFile(board, data)
		}
	}
}

func (blcu *BLCU) NotifyAck() {
	select {
	case blcu.ackChannel <- struct{}{}:
	default:
	}
}
