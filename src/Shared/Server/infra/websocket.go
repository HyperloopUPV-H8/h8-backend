package infra

import (
	"log"
	"net/http"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/server/infra/interfaces"
	"github.com/gorilla/websocket"
)

type Packet string

func (server HTTPServer[D, O, M]) HandleWebSocketOrder(route string, handler func(interfaces.WebSocket, chan<- O)) {
	server.router.Handle(route, SocketHandle[chan<- O]{
		function: handler,
		channel:  server.OrderChan,
	})
}

func (server HTTPServer[D, O, M]) HandleWebSocketData(route string, handler func(interfaces.WebSocket, <-chan D)) {
	server.router.Handle(route, SocketHandle[<-chan D]{
		function: handler,
		channel:  server.PacketChan,
	})
}

func (server HTTPServer[D, O, M]) HandleWebSocketMessage(route string, handler func(interfaces.WebSocket, <-chan M)) {
	server.router.Handle(route, SocketHandle[<-chan M]{
		function: handler,
		channel:  server.MessageChan,
	})
}

type SocketHandle[T any] struct {
	channel  T
	function func(interfaces.WebSocket, T)
}

var upgrader = websocket.Upgrader{}

func (handle SocketHandle[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("websocket handle: %s\n", err)
	}

	handle.function(conn, handle.channel)
}
