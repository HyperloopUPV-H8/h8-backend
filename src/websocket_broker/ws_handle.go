package websocket_broker

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/websocket_broker/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

type WebSocketBroker struct {
	handlers   map[string][]models.MessageHandler
	handlersMx *sync.Mutex
	clients    map[string]*websocket.Conn
	clientsMx  *sync.Mutex
	CloseChan  chan string
	trace      zerolog.Logger
}

func New() WebSocketBroker {
	trace.Info().Msg("new websocket broker")
	return WebSocketBroker{
		handlers:   make(map[string][]models.MessageHandler),
		handlersMx: &sync.Mutex{},
		clients:    make(map[string]*websocket.Conn),
		clientsMx:  &sync.Mutex{},
		CloseChan:  make(chan string),

		trace: trace.With().Str("component", "webSocketBroker").Logger(),
	}
}

func (broker *WebSocketBroker) HandleConn(writter http.ResponseWriter, request *http.Request) {
	broker.trace.Debug().Msg("new conn")
	defer request.Body.Close()

	writter.Header().Set("Access-Control-Allow-Origin", "*")

	upgrader := websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(writter, request, writter.Header())
	if err != nil {
		broker.trace.Error().Stack().Err(err).Msg("")
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		broker.trace.Error().Stack().Err(err).Msg("")
		conn.Close()
		return
	}

	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()

	broker.trace.Info().Str("id", id.String()).Msg("new client")
	broker.clients[id.String()] = conn
	go broker.readMessages(id.String(), conn)
}

func (broker *WebSocketBroker) readMessages(client string, conn *websocket.Conn) {
	broker.trace.Debug().Str("id", client).Msg("read messages")
	for {
		var message models.Message
		if err := conn.ReadJSON(&message); err != nil {
			broker.trace.Error().Str("id", client).Stack().Err(err).Msg("")
			broker.clientsMx.Lock()
			defer broker.clientsMx.Unlock()
			broker.closeClient(client)
			return
		}
		broker.updateHandlers(message.Topic, message.Payload, client)

	}
}

func (broker *WebSocketBroker) updateHandlers(topic string, payload json.RawMessage, source string) {
	broker.handlersMx.Lock()
	defer broker.handlersMx.Unlock()

	broker.trace.Trace().Str("topic", topic).Str("source", source).Msg("update")
	for _, handler := range broker.handlers[topic] {
		handler.UpdateMessage(topic, payload, source)
	}
}

func (broker *WebSocketBroker) sendMessage(topic string, payload any, targets ...string) error {

	broker.trace.Trace().Str("topic", topic).Strs("targets", targets).Msg("send message")
	message, err := models.NewMessage(topic, payload)
	if err != nil {
		broker.trace.Error().Stack().Err(err).Msg("")
		return err
	}

	if len(targets) == 0 {
		targets = common.Keys(broker.clients)
	}

	return broker.broadcastMessage(message, targets...)
}

func (broker *WebSocketBroker) broadcastMessage(message models.Message, targets ...string) error {
	broker.trace.Trace().Str("topic", message.Topic).Strs("targets", targets).Msg("broadcast message")
	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	var err error = nil
	for _, target := range targets {
		conn, ok := broker.clients[target]
		if !ok {
			broker.trace.Warn().Str("target", target).Msg("target not found")
			continue
		}

		if writeErr := conn.WriteJSON(message); writeErr != nil {
			broker.trace.Error().Stack().Err(writeErr).Msg("")
			broker.closeClient(target)
			err = writeErr
		}
	}
	return err
}

func (broker *WebSocketBroker) RegisterHandle(handler models.MessageHandler, topics ...string) {
	broker.handlersMx.Lock()
	defer broker.handlersMx.Unlock()
	broker.trace.Debug().Strs("topics", topics).Str("handler", handler.HandlerName()).Msg("register handle")
	handler.SetSendMessage(broker.sendMessage)
	for _, topic := range topics {
		broker.handlers[topic] = append(broker.handlers[topic], handler)
	}
}

func (broker *WebSocketBroker) RemoveHandler(topic string, handlerName string) {
	broker.handlersMx.Lock()
	defer broker.handlersMx.Unlock()
	broker.trace.Debug().Str("topic", topic).Str("handler", handlerName).Msg("remove handler")
	for i, handler := range broker.handlers[topic] {
		if handler.HandlerName() == handlerName {
			broker.handlers[topic][i] = broker.handlers[topic][len(broker.handlers[topic])-1]
			broker.handlers[topic] = broker.handlers[topic][:len(broker.handlers[topic])-1]
			return
		}
	}
}

func (broker *WebSocketBroker) closeClient(id string) error {
	broker.trace.Info().Str("id", id).Msg("close client")

	client, ok := broker.clients[id]

	if !ok {
		broker.trace.Warn().Str("target", id).Msg("client not found")
		return nil
	}

	if err := client.Close(); err != nil {
		broker.trace.Error().Stack().Err(err).Msg("")
		return err
	}
	delete(broker.clients, id)
	broker.CloseChan <- id
	return nil
}

func (broker *WebSocketBroker) Close() error {
	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	broker.trace.Info().Msg("close")
	var err error
	for client, conn := range broker.clients {
		if closeErr := conn.Close(); closeErr != nil {
			broker.trace.Error().Stack().Err(closeErr).Msg("")
			err = closeErr
			continue
		}
		delete(broker.clients, client)
	}
	return err
}
