package infra

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Packet string

func (server HTTPServer[R, S]) HandleWebSocketSend(route string, handler func(*websocket.Conn, chan S)) {
	server.router.Handle(route, SocketHandle[S]{
		function: handler,
		channel:  server.OrderSend,
	})
}

func (server HTTPServer[R, S]) HandleWebSocketRecv(route string, handler func(*websocket.Conn, chan R)) {
	server.router.Handle(route, SocketHandle[R]{
		function: handler,
		channel:  server.PacketRecv,
	})
}

type SocketHandle[T any] struct {
	channel  chan T
	function func(*websocket.Conn, chan T)
}

var upgrader = websocket.Upgrader{}

func (handle SocketHandle[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("websocket handle: %s\n", err)
	}

	handle.function(conn, handle.channel)

}
