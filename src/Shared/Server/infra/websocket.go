package infra

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Packet string

func (server HTTPServer[D, O, M]) HandleWebSocketOrder(route string, handler func(*websocket.Conn, chan O)) {
	server.router.Handle(route, SocketHandle[O]{
		function: handler,
		channel:  server.OrderSend,
	})
}

func (server HTTPServer[D, O, M]) HandleWebSocketData(route string, handler func(*websocket.Conn, chan D)) {
	server.router.Handle(route, SocketHandle[D]{
		function: handler,
		channel:  server.PacketRecv,
	})
}

func (server HTTPServer[D, O, M]) HandleWebSocketMessage(route string, handler func(*websocket.Conn, chan M)) {
	server.router.Handle(route, SocketHandle[M]{
		function: handler,
		channel:  server.MessageRecv,
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
