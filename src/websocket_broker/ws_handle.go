package websocket_broker

import (
	"fmt"
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
	clients    map[string]*websocket.Conn
	trace      zerolog.Logger
}

func New() WebSocketBroker {
	trace.Info().Msg("new websocket broker")
	return WebSocketBroker{
		handlersMx: &sync.Mutex{},
		handlers:   make(map[string][]models.MessageHandler),
		clientsMx:  &sync.Mutex{},
		clients:    make(map[string]*websocket.Conn),
		trace:      trace.With().Str("component", "webSocketBroker").Logger(),
	}
}

func (broker *WebSocketBroker) Add(conn *websocket.Conn) error {
	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	id := conn.RemoteAddr().String()
	broker.clients[id] = conn
	go broker.readMessages(id, conn)
	go broker.ping(id, conn)
	return nil
}

func (broker *WebSocketBroker) readMessages(id string, conn *websocket.Conn) {
	broker.trace.Debug().Str("id", conn.RemoteAddr().String()).Msg("read messages")
	for {
		var msg models.Message
		if err := conn.ReadJSON(&msg); err != nil {
			broker.removeClient(id)
			return
		}
		broker.updateHandlers(id, msg)
	}
}

func (broker *WebSocketBroker) ping(id string, conn *websocket.Conn) {
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		broker.clientsMx.Lock()
		if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			broker.unsafeRemoveClient(id)
			broker.clientsMx.Unlock()
			return
		}
		broker.clientsMx.Unlock()
	}
}

func (broker *WebSocketBroker) updateHandlers(clientId string, msg models.Message) {
	broker.handlersMx.Lock()
	defer broker.handlersMx.Unlock()

	broker.trace.Trace().Str("topic", msg.Topic).Str("clientId", clientId).Msg("update")
	for _, handler := range broker.handlers[msg.Topic] {
		handler.UpdateMessage(msg.Topic, msg.Payload, clientId)
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
		return broker.broadcastMessage(message)
	}

	return broker.sendToTargets(message, targets)
}

func (broker *WebSocketBroker) sendToTargets(msg models.Message, targets []string) error {
	failedTargets := make([]string, 0)
	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()

	for _, target := range targets {
		conn, ok := broker.clients[target]

		if !ok {
			broker.trace.Warn().Str("target", target).Msg("target not found")
			failedTargets = append(failedTargets, target)
			continue
		}

		if err := conn.WriteJSON(msg); err != nil {
			broker.trace.Error().Stack().Err(err).Msg("")
			delete(broker.clients, target)
			failedTargets = append(failedTargets, target)
		}
	}

	if len(failedTargets) == 0 {
		return nil
	}

	return fmt.Errorf("failed targets: %v", failedTargets)
}

func (broker *WebSocketBroker) broadcastMessage(message models.Message) error {
	broker.trace.Trace().Str("topic", message.Topic).Msg("broadcast message")
	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	var err error
	for id, conn := range broker.clients {
		if err = conn.WriteJSON(message); err != nil {
			delete(broker.clients, id)
			broker.trace.Error().Stack().Err(err).Msg("")
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

func (broker *WebSocketBroker) Close() error {
	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	broker.trace.Info().Msg("close")
	var err error
	for id, conn := range broker.clients {
		if closeErr := conn.Close(); closeErr != nil {
			broker.trace.Error().Stack().Err(closeErr).Msg("")
			err = closeErr
			continue
		}
		delete(broker.clients, id)
	}

	return err
}

func (broker *WebSocketBroker) removeClient(id string) {
	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	broker.unsafeRemoveClient(id)

}

func (broker *WebSocketBroker) unsafeRemoveClient(id string) {
	_, ok := broker.clients[id]

	if !ok {
		broker.trace.Error().Str("clientId", id).Msg("client to be removed not found")
		return
	}

	delete(broker.clients, id)
	broker.trace.Error().Str("clientId", id).Msg("remove client")
}
