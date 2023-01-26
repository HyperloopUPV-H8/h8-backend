package websocket_handle

import (
	"log"
	"net/http"
	"sync"

	"github.com/HyperloopUPV-H8/Backend-H8/websocket_handle/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WSHandle struct {
	handles   map[string]chan models.MessageTarget
	clients   map[string]chan models.Message
	clientsMx sync.Mutex
}

func RunWSHandle(router *mux.Router, route string, handles map[string]chan models.MessageTarget) *WSHandle {
	handle := &WSHandle{
		handles:   handles,
		clients:   make(map[string]chan models.Message),
		clientsMx: sync.Mutex{},
	}

	router.HandleFunc(route, handle.handleConn)

	go handle.runRecv()
	go handle.runSend()

	return handle
}

func (handle *WSHandle) multiplex(source string, msg models.Message) {
	handle.handles[msg.Type] <- models.MessageTarget{
		Target: []string{source},
		Msg:    msg,
	}
}

func (handle *WSHandle) runRecv() {
	for {
		for client, messages := range handle.clients {
			select {
			case msg := <-messages:
				handle.multiplex(client, msg)
			default:
			}
		}
	}
}

func (handle *WSHandle) distribute(msg models.MessageTarget) {
	handle.clientsMx.Lock()
	defer handle.clientsMx.Unlock()
	for _, target := range msg.Target {
		handle.clients[target] <- msg.Msg
	}

	if len(msg.Target) == 0 {
		for _, client := range handle.clients {
			client <- msg.Msg
		}
	}
}

func (handle *WSHandle) runSend() {
	for {
		for _, messages := range handle.handles {
			select {
			case msg := <-messages:
				handle.distribute(msg)
			default:
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (handle *WSHandle) handleConn(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebsocketHandler: handleConn: %s\n", err)
		return
	}

	handle.clientsMx.Lock()
	defer handle.clientsMx.Unlock()
	handle.clients[conn.RemoteAddr().String()] = handleSocket(conn)
}

func handleSocket(conn *websocket.Conn) chan models.Message {
	messages := make(chan models.Message)

	go func(conn *websocket.Conn, messages chan<- models.Message) {
		for {
			msg := new(models.Message)
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("WebsocketHandle: handleSocket: %s\n", err)
				conn.Close()
				return
			}
			messages <- *msg
		}
	}(conn, messages)

	go func(conn *websocket.Conn, messages <-chan models.Message) {
		for msg := range messages {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Printf("WebsocketHandle: handleSocket: %s\n", err)
				conn.Close()
				return
			}
		}
	}(conn, messages)

	return messages
}
