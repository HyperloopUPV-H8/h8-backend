package server

import (
	"net/http"

	"github.com/gorilla/websocket"
	trace "github.com/rs/zerolog/log"
)

type ConnectionHandler interface {
	Add(*websocket.Conn) error
}

func (server *WebServer) serveWebsocket(path string, upgrader *websocket.Upgrader, headers map[string]string) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		trace.Info().Str("server", server.name).Msg("websocket connection")
		for key, value := range headers {
			w.Header().Set(key, value)
		}

		if server.config.MaxConnections != nil && server.connected.Load() >= *server.config.MaxConnections {
			http.Error(w, "Max connections reached", http.StatusTooManyRequests)
			return
		}

		conn, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			return
		}

		conn.SetCloseHandler(func(code int, text string) error {
			server.connected.Add(-1)
			return nil
		})

		err = server.connHandler.Add(conn)
		if err != nil {
			return
		}
		server.connected.Add(1)

	}

	server.router.HandleFunc(path, handler)
}
