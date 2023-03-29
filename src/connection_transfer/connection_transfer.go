package connection_transfer

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/HyperloopUPV-H8/Backend-H8/connection_transfer/models"
)

const (
	CONNECTION_TRANSFER_HANDLER_NAME = "connectionTransfer"
	CONNECTION_TRANSFER_TOPIC        = "connection/update"
)

var (
	connectionTransfer *ConnectionTransfer
)

func Get() *ConnectionTransfer {
	if connectionTransfer == nil {
		initConnectionTransfer()
	}
	return connectionTransfer
}

func initConnectionTransfer() {
	connectionTransfer = &ConnectionTransfer{
		writeMx:     &sync.Mutex{},
		boardStatus: make(map[string]models.Connection),
		sendMessage: defaultSendMessage,
	}
}

type ConnectionTransfer struct {
	writeMx     *sync.Mutex
	boardStatus map[string]models.Connection
	sendMessage func(topic string, payload any, target ...string) error
}

func (connectionTransfer *ConnectionTransfer) UpdateMessage(topic string, payload json.RawMessage, source string) {
	connectionTransfer.send()
}

func (connectionTransfer *ConnectionTransfer) SetSendMessage(sendMessage func(topic string, payload any, target ...string) error) {
	connectionTransfer.sendMessage = sendMessage
}

func (connectionTransfer *ConnectionTransfer) HandlerName() string {
	return CONNECTION_TRANSFER_HANDLER_NAME
}

func (connectionTransfer *ConnectionTransfer) Update(name string, up bool) {
	connectionTransfer.writeMx.Lock()
	defer connectionTransfer.writeMx.Unlock()
	connectionTransfer.boardStatus[name] = models.Connection{
		Name:        name,
		IsConnected: up,
	}

	connectionTransfer.send()
}

func (connectionTransfer *ConnectionTransfer) send() {
	if err := connectionTransfer.sendMessage(CONNECTION_TRANSFER_TOPIC, connectionTransfer.boardStatus); err != nil {
		log.Printf("ConnectionTransfer: send: sendMessage: %s\n", err)
	}
}

func defaultSendMessage(string, any, ...string) error {
	return errors.New("connection transfer must be registered before use")
}
