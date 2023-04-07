package blcu

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	excel_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

// TODO: Get these values from the TOML
const (
	BLCU_COMPONENT_NAME       = "blcu"
	BLCU_BOARD_NAME           = "blcu"
	BLCU_HANDLER_NAME         = "blcu"
	BLCU_UPLOAD_ORDER_ID      = 700
	BLCU_DOWNLOAD_ORDER_ID    = 701
	BLCU_UPLOAD_ORDER_FIELD   = "board"
	BLCU_DOWNLOAD_ORDER_FIELD = "board"
	BLCU_INPUT_CHAN_BUF       = 100
	BLCU_ACK_CHAN_BUF         = 1
	BLCU_ACK_PACKET_NAME      = "tftp_ack"
)

type BLCU struct {
	addr  string
	ackID uint16

	inputChannel chan models.Update
	ackChannel   chan struct{}

	sendOrder   func(models.Order) error
	sendMessage func(topic string, payload any, targets ...string) error

	trace zerolog.Logger
}

func NewBLCU() *BLCU {
	trace.Info().Msg("New BLCU")
	blcu := &BLCU{
		inputChannel: make(chan models.Update, BLCU_INPUT_CHAN_BUF),
		ackChannel:   make(chan struct{}, BLCU_ACK_CHAN_BUF),
		trace:        trace.With().Str("component", BLCU_COMPONENT_NAME).Logger(),
	}

	return blcu
}

func (blcu *BLCU) AddGlobal(global excel_models.GlobalInfo) {
	blcu.addr = fmt.Sprintf("%s:%s", global.BoardToIP["BLCU"], global.ProtocolToPort["TFTP"])
	blcu.trace.Debug().Str("addr", blcu.addr).Msg("Add global info")
}

func (blcu *BLCU) AddPacket(boardName string, packet excel_models.Packet) {
	if packet.Description.Name != BLCU_ACK_PACKET_NAME {
		return
	}

	id, err := strconv.ParseUint(packet.Description.ID, 10, 16)
	if err != nil {
		blcu.trace.Error().Err(err).Stack().Msg("Error parsing packet ID")
		return
	}

	blcu.ackID = uint16(id)
	blcu.trace.Debug().Uint16("ackID", blcu.ackID).Msg("Add packet info")
}

func (blcu *BLCU) UpdateMessage(topic string, payload json.RawMessage, source string) {
	blcu.trace.Debug().Str("topic", topic).Str("source", source).Msg("Update message")
	switch topic {
	case os.Getenv("BLCU_UPLOAD_TOPIC"):
		if err := blcu.handleUpload(payload); err != nil {
			blcu.notifyUploadFailure()
		} else {
			blcu.notifyUploadSuccess()
		}
	case os.Getenv("BLCU_DOWNLOAD_TOPIC"):
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
		select {
		case blcu.ackChannel <- struct{}{}:
		default:
		}
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
