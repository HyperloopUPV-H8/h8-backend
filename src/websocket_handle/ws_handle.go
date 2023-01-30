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
	handles map[string]chan models.MessageTarget
	conns   map[string]*websocket.Conn
	connsMx sync.Mutex
}

func RunWSHandle(router *mux.Router, route string, handles map[string]chan models.MessageTarget) *WSHandle {
	handle := &WSHandle{
		handles: handles,
		conns:   make(map[string]*websocket.Conn),
		connsMx: sync.Mutex{},
	}

	router.HandleFunc(route, handle.handleConn)

	go handle.handleTx()

	return handle
}

func (handle *WSHandle) handleTx() {
	for {
		for _, output := range handle.handles {
			select {
			case msg := <-output:
				if len(msg.Target) == 0 {
					for _, conn := range handle.conns {
						err := conn.WriteJSON(msg.Msg)
						if err != nil {
							handle.Close(conn)
						}
					}
				} else {
					for _, target := range msg.Target {
						err := handle.conns[target].WriteJSON(msg.Msg)
						if err != nil {
							handle.Close(handle.conns[target])
						}
					}
				}
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
		log.Printf("WebSocketHandler: handleConn: %s\n", err)
		return
	}

	go handle.handleSocket(conn)
	handle.conns[conn.RemoteAddr().String()] = conn
}

func (handle *WSHandle) handleSocket(conn *websocket.Conn) {
	defer handle.Close(conn)
	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			return
		}
		handle.handles[msg.Topic] <- models.MessageTarget{
			Target: []string{conn.RemoteAddr().String()},
			Msg:    msg,
		}
	}

}

func (handle *WSHandle) Close(conn *websocket.Conn) {
	handle.connsMx.Lock()
	defer handle.connsMx.Unlock()
	delete(handle.conns, conn.RemoteAddr().String())
	conn.Close()
}
