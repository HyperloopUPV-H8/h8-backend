package websocket_broker

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/websocket_broker/models"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

type WebSocketBroker struct {
	handlersMx *sync.Mutex
	handlers   map[string][]models.MessageHandler
	clientsMx  *sync.Mutex
	clients    map[string]map[string]*websocket.Conn
	removed    map[string]chan struct{}
	CloseChan  chan string
	trace      zerolog.Logger
}

func New() WebSocketBroker {
	trace.Info().Msg("new websocket broker")
	return WebSocketBroker{
		handlersMx: &sync.Mutex{},
		handlers:   make(map[string][]models.MessageHandler),
		clientsMx:  &sync.Mutex{},
		clients:    make(map[string]map[string]*websocket.Conn),
		removed:    make(map[string]chan struct{}),
		CloseChan:  make(chan string, 100),
		trace:      trace.With().Str("component", "webSocketBroker").Logger(),
	}
}

func (broker *WebSocketBroker) AddServer(server string) {
	broker.trace.Debug().Str("server", server).Msg("add server")
	broker.handlersMx.Lock()
	defer broker.handlersMx.Unlock()
	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	broker.handlers[server] = make([]models.MessageHandler, 0)
	broker.clients[server] = make(map[string]*websocket.Conn)
	broker.removed[server] = make(chan struct{}, 1)
}

func (broker *WebSocketBroker) Add(server string, conn *websocket.Conn) error {
	broker.trace.Debug().Str("server", server).Msg("add")
	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	if _, ok := broker.clients[server]; !ok {
		broker.clients[server] = make(map[string]*websocket.Conn)
	}
	broker.clients[server][conn.RemoteAddr().String()] = conn
	go broker.readMessages(server, conn.RemoteAddr().String(), conn)
	go broker.ping(server, conn.RemoteAddr().String(), conn)
	return nil
}

func (broker *WebSocketBroker) Remove(server string) <-chan struct{} {
	return broker.removed[server]
}

func (broker *WebSocketBroker) closeClientBlocking(server string, client string) {
	broker.trace.Debug().Str("id", client).Msg("close client blocking")
	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	broker.closeClient(server, client)
}

func (broker *WebSocketBroker) readMessages(server string, client string, conn *websocket.Conn) {
	broker.trace.Debug().Str("id", client).Msg("read messages")
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		var message models.Message
		if err := conn.ReadJSON(&message); err != nil {
			broker.closeClientBlocking(server, client)
			return
		}
		broker.updateHandlers(message.Topic, message.Payload, client)

	}
}

func (broker *WebSocketBroker) ping(server string, client string, conn *websocket.Conn) {
	broker.trace.Debug().Str("id", client).Msg("ping")
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			broker.closeClientBlocking(server, client)
			return
		}
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
		targets = broker.getAllTargets()
	}

	return broker.broadcastServers(message, targets...)
}

func (broker *WebSocketBroker) getAllTargets() []string {
	targets := make([]string, 0)
	for _, server := range broker.clients {
		for target := range server {
			targets = append(targets, target)
		}
	}
	return targets
}

func (broker *WebSocketBroker) broadcastServers(message models.Message, targets ...string) error {
	var err error
	for server := range broker.clients {
		writeErr := broker.broadcastMessage(server, message, targets...)
		if writeErr != nil {
			err = writeErr
		}
	}
	return err
}

func (broker *WebSocketBroker) broadcastMessage(server string, message models.Message, targets ...string) error {
	broker.trace.Trace().Str("topic", message.Topic).Strs("targets", targets).Msg("broadcast message")
	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	var err error = nil
	for _, target := range targets {
		conn, ok := broker.clients[server][target]
		if !ok {
			broker.trace.Warn().Str("target", target).Msg("target not found")
			continue
		}

		if writeErr := conn.WriteJSON(message); writeErr != nil {
			broker.trace.Error().Stack().Err(writeErr).Msg("")
			broker.closeClient(server, target)
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

func (broker *WebSocketBroker) closeClient(server string, id string) error {
	broker.trace.Info().Str("id", id).Msg("close client")

	client, ok := broker.clients[server][id]

	if !ok {
		broker.trace.Warn().Str("target", id).Msg("client not found")
		return nil
	}

	if err := client.Close(); err != nil {
		broker.trace.Error().Stack().Err(err).Msg("")
		return err
	}
	delete(broker.clients, id)
	broker.removed[server] <- struct{}{}
	broker.CloseChan <- id
	return nil
}

func (broker *WebSocketBroker) Close() error {
	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	broker.trace.Info().Msg("close")
	var err error
	for _, server := range broker.clients {
		for client, conn := range server {
			if closeErr := conn.Close(); closeErr != nil {
				broker.trace.Error().Stack().Err(closeErr).Msg("")
				err = closeErr
				continue
			}
			delete(broker.clients, client)
		}
	}

	return err
}
