package blcu

import (
	"encoding/json"
	"fmt"
	"strconv"

	excel_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

type BLCU struct {
	addr       string
	boardToId  map[string]uint16
	ackChannel chan struct{}

	sendOrder   func(models.Order) error
	sendMessage func(topic string, payload any, targets ...string) error

	config BLCUConfig

	trace zerolog.Logger
}

func NewBLCU(global excel_models.GlobalInfo, config BLCUConfig) BLCU {
	trace.Info().Msg("New BLCU")
	blcu := BLCU{
		addr:        fmt.Sprintf("%s:%s", global.BoardToIP["BLCU"], global.ProtocolToPort["TFTP"]),
		boardToId:   getBoardToId(global.BoardToId),
		ackChannel:  make(chan struct{}, BLCU_ACK_CHAN_BUF),
		trace:       trace.With().Str("component", BLCU_COMPONENT_NAME).Logger(),
		config:      config,
		sendOrder:   func(o models.Order) error { return nil },
		sendMessage: func(topic string, payload any, targets ...string) error { return nil },
	}

	return blcu
}

func getBoardToId(boardToIdStr map[string]string) map[string]uint16 {
	boardToId := make(map[string]uint16)

	for name, idStr := range boardToIdStr {
		id, err := strconv.Atoi(idStr)

		if err != nil {
			trace.Error().Err(err).Stack().Msg("Convert board id")
			continue
		}

		boardToId[name] = uint16(id)
	}

	return boardToId
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
		result := blcu.upload(payload)
		if result.Err != nil {
			blcu.notifyUploadFailure()
		} else {
			blcu.notifyUploadSuccess()
		}

	case blcu.config.Topics.Download:
		board, result := blcu.download(payload)
		if result.Err != nil {
			blcu.notifyDownloadFailure()
		} else {
			blcu.notifyDownloadSuccess(result.Data)
			blcu.writeDownloadFile(board, result.Data)
		}
	}
}

func (blcu *BLCU) NotifyAck() {
	select {
	case blcu.ackChannel <- struct{}{}:
	default:
	}
}
