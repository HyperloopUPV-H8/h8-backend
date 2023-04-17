package blcu

import (
	"encoding/json"
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	excel_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

type BLCU struct {
	addr  string
	ackID uint16

	inputChannel chan models.Update
	ackChannel   chan struct{}

	sendOrder   func(models.Order) error
	sendMessage func(topic string, payload any, targets ...string) error

	config BLCUConfig

	trace zerolog.Logger
}

func NewBLCU(global excel_models.GlobalInfo, config BLCUConfig) BLCU {
	trace.Info().Msg("New BLCU")
	blcu := BLCU{
		addr:         fmt.Sprintf("%s:%s", global.BoardToIP["BLCU"], global.ProtocolToPort["TFTP"]),
		inputChannel: make(chan models.Update, BLCU_INPUT_CHAN_BUF),
		ackChannel:   make(chan struct{}, BLCU_ACK_CHAN_BUF),
		trace:        trace.With().Str("component", BLCU_COMPONENT_NAME).Logger(),
		config:       config,
	}

	return blcu
}

func (blcu *BLCU) UpdateMessage(topic string, payload json.RawMessage, source string) {
	blcu.trace.Debug().Str("topic", topic).Str("source", source).Msg("Update message")
	switch topic {
	case blcu.config.Topics.Upload:
		if err := blcu.handleUpload(payload); err != nil {
			blcu.notifyUploadFailure()
		} else {
			blcu.notifyUploadSuccess()
		}
	case blcu.config.Topics.Download:
		if file, err := blcu.handleDownload(payload); err != nil {
			blcu.notifyDownloadFailure()
		} else {
			blcu.notifyDownloadSuccess(file)
		}
	}
}

func (blcu *BLCU) SetSendMessage(sendMessage func(topic string, payload any, targets ...string) error) {
	blcu.trace.Debug().Msg("Set send message")
	blcu.sendMessage = sendMessage
}

func (blcu *BLCU) HandlerName() string {
	return BLCU_HANDLER_NAME
}

func (blcu *BLCU) Request(order models.Order) error {
	return blcu.sendOrder(order)
}

func (blcu *BLCU) Listen(destination chan<- models.Update) {
	blcu.trace.Debug().Msg("Listen")
	for update := range blcu.inputChannel {
		destination <- update
	}
}

func (blcu *BLCU) Input(update models.Update) {
	if update.ID == blcu.ackID {
		common.TrySend(blcu.ackChannel, struct{}{})
	}
	blcu.inputChannel <- update
}

func (blcu *BLCU) Output(sendOrder func(models.Order) error) {
	blcu.trace.Debug().Msg("Set output")
	blcu.sendOrder = sendOrder
}

func (blcu *BLCU) Name() string {
	return BLCU_BOARD_NAME
}
