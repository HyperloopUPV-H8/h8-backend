package websocket_broker

import (
	"fmt"
	"sync"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/websocket_broker/models"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type WebSocketBroker struct {
	handlersMx *sync.Mutex
	handlers   map[string][]models.MessageHandler
	clientsMx  *sync.Mutex
	clients    map[string]*websocket.Conn
	log        zerolog.Logger
}

func New() WebSocketBroker {
	log.Info().Msg("new websocket broker")
	return WebSocketBroker{
		handlersMx: &sync.Mutex{},
		handlers:   make(map[string][]models.MessageHandler),
		clientsMx:  &sync.Mutex{},
		clients:    make(map[string]*websocket.Conn),
		log:        log.With().Str("component", "webSocketBroker").Logger(),
	}
}

func (broker *WebSocketBroker) Add(conn *websocket.Conn) error {
	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	id := conn.RemoteAddr().String()
	broker.clients[id] = conn
	broker.log.Info().Str("id", id).Msg("add client")
	go broker.readMessages(id, conn)
	go broker.ping(id, conn)
	return nil
}

func (broker *WebSocketBroker) readMessages(id string, conn *websocket.Conn) {
	broker.log.Debug().Str("id", id).Msg("read messages")
	defer broker.log.Debug().Str("id", id).Msg("stop reading messages")
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
	broker.log.Debug().Str("id", id).Msg("send pings")
	defer broker.log.Debug().Str("id", id).Msg("stop sending pings")
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		if broker.sendPing(id, conn) != nil {
			return
		}
	}
}

func (broker *WebSocketBroker) sendPing(id string, conn *websocket.Conn) error {
	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	broker.log.Trace().Str("id", id).Msg("ping")
	if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
		broker.unsafeRemoveClient(id)
		return err
	}
	return nil
}

func (broker *WebSocketBroker) updateHandlers(id string, msg models.Message) {
	broker.handlersMx.Lock()
	defer broker.handlersMx.Unlock()

	broker.log.Trace().Str("topic", msg.Topic).Str("id", id).Msg("update")
	for _, handler := range broker.handlers[msg.Topic] {
		handler.UpdateMessage(msg.Topic, msg.Payload, id)
	}
}

func (broker *WebSocketBroker) sendMessage(topic string, payload any, targets ...string) error {
	broker.log.Trace().Str("topic", topic).Strs("targets", targets).Msg("send message")
	message, err := models.NewMessage(topic, payload)
	if err != nil {
		broker.log.Error().Stack().Err(err).Str("topic", topic).Any("payload", payload).Msg("create message")
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
			broker.log.Debug().Str("target", target).Msg("target not found")
			failedTargets = append(failedTargets, target)
			continue
		}

		broker.unsafeWriteToClient(target, conn, msg)
	}

	if len(failedTargets) == 0 {
		return nil
	}

	return fmt.Errorf("failed targets: %v", failedTargets)
}

func (broker *WebSocketBroker) broadcastMessage(msg models.Message) error {
	broker.log.Trace().Str("topic", msg.Topic).Msg("broadcast message")
	broker.clientsMx.Lock()
	defer broker.clientsMx.Unlock()
	var err error
	for id, conn := range broker.clients {
		broker.unsafeWriteToClient(id, conn, msg)
	}

	return err
}

func (broker *WebSocketBroker) unsafeWriteToClient(id string, conn *websocket.Conn, msg models.Message) {
	broker.log.Trace().Str("id", id).Str("topic", msg.Topic).Msg("send message")
	if err := conn.WriteJSON(msg); err != nil {
		broker.log.Error().Stack().Err(err).Str("id", id).Msg("write to client")
		broker.unsafeRemoveClient(id)
	}
}

func (broker *WebSocketBroker) RegisterHandle(handler models.MessageHandler, topics ...string) {
	broker.handlersMx.Lock()
	defer broker.handlersMx.Unlock()
	broker.log.Debug().Strs("topics", topics).Str("handler", handler.HandlerName()).Msg("register handle")
	handler.SetSendMessage(broker.sendMessage)
	for _, topic := range topics {
		broker.handlers[topic] = append(broker.handlers[topic], handler)
	}
}

func (broker *WebSocketBroker) RemoveHandler(topic string, handlerName string) {
	broker.handlersMx.Lock()
	defer broker.handlersMx.Unlock()
	broker.log.Debug().Str("topic", topic).Str("handler", handlerName).Msg("remove handler")
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
	broker.log.Info().Msg("close")
	var err error
	for id, conn := range broker.clients {
		broker.log.Debug().Str("id", id).Msg("close client")
		if closeErr := conn.Close(); closeErr != nil {
			broker.log.Error().Stack().Err(closeErr).Str("id", id).Msg("close client")
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
	broker.log.Info().Str("id", id).Msg("remove client")
	_, ok := broker.clients[id]

	if !ok {
		broker.log.Debug().Str("id", id).Msg("client to be removed not found")
		return
	}

	delete(broker.clients, id)
}
