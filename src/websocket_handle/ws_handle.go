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
		var delete []*websocket.Conn
		handle.connsMx.Lock()
		for _, output := range handle.handles {
			select {
			case msg := <-output:
				if len(msg.Target) == 0 {
					delete = handle.broadcast(msg)
				} else {
					delete = handle.unicast(msg)
				}
			default:
			}
		}
		handle.connsMx.Unlock()
		for _, conn := range delete {
			handle.Close(conn)
		}
	}
}

func (handle *WSHandle) broadcast(msg models.MessageTarget) []*websocket.Conn {
	delete := make([]*websocket.Conn, 0, len(handle.handles))
	for _, conn := range handle.conns {
		err := conn.WriteJSON(msg.Msg)
		if err != nil {
			delete = append(delete, conn)
		}
	}
	return delete
}

func (handle *WSHandle) unicast(msg models.MessageTarget) []*websocket.Conn {
	delete := make([]*websocket.Conn, 0, len(handle.handles))
	for _, target := range msg.Target {
		err := handle.conns[target].WriteJSON(msg.Msg)
		if err != nil {
			delete = append(delete, handle.conns[target])
		}
	}
	return delete
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
	handle.connsMx.Lock()
	handle.conns[conn.RemoteAddr().String()] = conn
	handle.connsMx.Unlock()
}

func (handle *WSHandle) handleSocket(conn *websocket.Conn) {
	defer handle.Close(conn)
	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			return
		}
		handle.connsMx.Lock()
		handle.handles[msg.Topic] <- models.MessageTarget{
			Target: []string{conn.RemoteAddr().String()},
			Msg:    msg,
		}
		handle.connsMx.Unlock()
	}

}

func (handle *WSHandle) Close(conn *websocket.Conn) {
	handle.connsMx.Lock()
	defer handle.connsMx.Unlock()
	log.Printf("closed %s\n", conn.RemoteAddr().String())
	delete(handle.conns, conn.RemoteAddr().String())
	delete(handle.handles, conn.RemoteAddr().String())
	conn.Close()
}
