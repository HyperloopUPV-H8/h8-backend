package data_transfer

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

const (
	DATA_TRANSFER_HANDLER_NAME = "dataTransfer"
	DATA_TRANSFER_REFRESH_RATE = time.Second / 30
	DATA_TRANSFER_TOPIC        = "packet/update"
)

var (
	dataTransfer *DataTransfer
)

func Get() *DataTransfer {
	if dataTransfer == nil {
		initDataTransfer()
	}
	return dataTransfer
}

func initDataTransfer() {
	dataTransfer = &DataTransfer{
		bufMx:       &sync.Mutex{},
		packetBuf:   make(map[uint16]models.Update),
		ticker:      time.NewTicker(DATA_TRANSFER_REFRESH_RATE),
		sendMessage: defaultSendMessage,
	}

	go dataTransfer.run()
}

type DataTransfer struct {
	bufMx       *sync.Mutex
	packetBuf   map[uint16]models.Update
	ticker      *time.Ticker
	sendMessage func(topic string, payload any, target ...string) error
}

func (dataTransfer *DataTransfer) UpdateMessage(string, json.RawMessage, string) {}

func (dataTransfer *DataTransfer) SetSendMessage(sendMessage func(topic string, payload any, target ...string) error) {
	dataTransfer.sendMessage = sendMessage
}

func (DataTransfer *DataTransfer) HandlerName() string {
	return DATA_TRANSFER_HANDLER_NAME
}

func (dataTransfer *DataTransfer) run() {
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
	if err := dataTransfer.sendMessage(DATA_TRANSFER_TOPIC, dataTransfer.packetBuf); err != nil {
		log.Printf("DataTransfer: sendBuf: sendMessage: %s\n", err)
		return
	}
	dataTransfer.packetBuf = make(map[uint16]models.Update, len(dataTransfer.packetBuf))
}

func (dataTransfer *DataTransfer) Update(update models.Update) {
	dataTransfer.bufMx.Lock()
	defer dataTransfer.bufMx.Unlock()
	dataTransfer.packetBuf[update.ID] = update
}

func defaultSendMessage(string, any, ...string) error {
	return errors.New("data transfer must be registered before use")
}
