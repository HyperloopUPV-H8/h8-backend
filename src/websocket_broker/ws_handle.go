package websocket_broker

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/websocket_broker/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	upgrader websocket.Upgrader = websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool { return true },
	}
	broker *WebSocketBroker = nil
)

func Get() *WebSocketBroker {
	if broker == nil {
		initBroker()
	}
	return broker
}

func initBroker() {
	broker = &WebSocketBroker{
		handlers:   make(map[string][]models.MessageHandler),
		handlersMx: &sync.Mutex{},
		clients:    make(map[string]*websocket.Conn),
		clientsMx:  &sync.Mutex{},
	}
}

type WebSocketBroker struct {
	handlers   map[string][]models.MessageHandler
	handlersMx *sync.Mutex
	clients    map[string]*websocket.Conn
	clientsMx  *sync.Mutex
}

func (broker *WebSocketBroker) HandleConn(writter http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	writter.Header().Set("Access-Control-Allow-Origin", "*")

	conn, err := upgrader.Upgrade(writter, request, writter.Header())
	if err != nil {
		log.Printf("WebSocketBroker: handleConn: Upgrade: %s\n", err)
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		log.Printf("WebSocketBroker: handleConn: uuid.NewRandom: %s\n", err)
		conn.Close()
		return
	}

	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	broker.clients[id.String()] = conn
	go broker.readMessages(id.String(), conn)
}

func (broker *WebSocketBroker) readMessages(client string, conn *websocket.Conn) {
	for {
		var message models.Message
		if err := conn.ReadJSON(&message); err != nil {
			log.Printf("WebSocketBroker: readMessages: ReadJSON: %s\n", err)
			broker.closeClient(client)
			return
		}

		broker.updateHandlers(message.Topic, message.Payload, client)
	}
}

func (broker *WebSocketBroker) updateHandlers(topic string, payload json.RawMessage, source string) {
	broker.handlersMx.Lock()
	defer broker.handlersMx.Unlock()
	for _, handler := range broker.handlers[topic] {
		handler.UpdateMessage(topic, payload, source)
	}
}

func (broker *WebSocketBroker) sendMessage(topic string, payload any, targets ...string) error {
	message, err := models.NewMessage(topic, payload)
	if err != nil {
		return err
	}

	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	if len(targets) == 0 {
		targets = common.Keys(broker.clients)
	}

	return broker.broadcastMessage(message, targets...)
}

func (broker *WebSocketBroker) broadcastMessage(message models.Message, targets ...string) error {
	var err error = nil
	for _, target := range targets {
		conn, ok := broker.clients[target]
		if !ok {
			continue
		}

		if writeErr := conn.WriteJSON(message); writeErr != nil {
			log.Printf("WebSocketBroker: broadcastMessage: WriteJSON: %s\n", writeErr)
			go broker.closeClient(target)
			err = writeErr
		}
	}
	return err
}

func (broker *WebSocketBroker) RegisterHandle(handler models.MessageHandler, topics ...string) {
	broker.handlersMx.Lock()
	defer broker.handlersMx.Unlock()
	handler.SetSendMessage(broker.sendMessage)
	for _, topic := range topics {
		broker.handlers[topic] = append(broker.handlers[topic], handler)
	}
}

func (broker *WebSocketBroker) RemoveHandler(topic string, handlerName string) {
	broker.handlersMx.Lock()
	defer broker.handlersMx.Unlock()
	for i, handler := range broker.handlers[topic] {
		if handler.HandlerName() == handlerName {
			broker.handlers[topic][i] = broker.handlers[topic][len(broker.handlers[topic])-1]
			broker.handlers[topic] = broker.handlers[topic][:len(broker.handlers[topic])-1]
			return
		}
	}
}

func (broker *WebSocketBroker) closeClient(id string) error {
	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	if err := broker.clients[id].Close(); err != nil {
		return err
	}
	delete(broker.clients, id)
	return nil
}

func (broker *WebSocketBroker) Close() error {
	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	var err error
	for client, conn := range broker.clients {
		if closeErr := conn.Close(); closeErr != nil {
			err = closeErr
			continue
		}
		delete(broker.clients, client)
	}
	return err
}
