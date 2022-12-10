package models

import "github.com/gorilla/websocket"

type WebsocketHandler interface {
	HandleConn(*websocket.Conn)
}
