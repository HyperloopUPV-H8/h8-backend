package server

import (
	"encoding/json"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WebServer struct {
	name        string
	router      *mux.Router
	connections ConnectionHandler
	connected   *atomic.Int32
	config      ServerConfig
}

func NewWebServer(name string, connectionHandle ConnectionHandler, staticData EndpointData, config ServerConfig) (*WebServer, error) {
	server := &WebServer{
		name:        name,
		router:      mux.NewRouter(),
		connections: connectionHandle,
		connected:   &atomic.Int32{},
		config:      config,
	}

	headers := map[string]string{
		"Access-Control-Allow-Origin": "*",
	}

	err := server.serveJSON("/backend"+config.Endpoints.PodData, staticData.PodData, headers)
	if err != nil {
		return nil, err
	}

	err = server.serveJSON("/backend"+config.Endpoints.OrderData, staticData.OrderData, headers)
	if err != nil {
		return nil, err
	}

	err = server.serveJSON("/backend"+config.Endpoints.ProgramableBoards, staticData.ProgramableBoards, headers)
	if err != nil {
		return nil, err
	}

	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	server.serveWebsocket(config.Endpoints.Connections, upgrader, headers)
	go server.consumeRemoved()

	server.serveFiles(config.Endpoints.Files, config.StaticPath)

	if err != nil {
		return nil, err
	}

	return server, nil
}

func (server *WebServer) serveJSON(path string, data any, headers map[string]string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		for key, value := range headers {
			w.Header().Set(key, value)
		}
		w.Write(jsonData)
	}

	server.router.HandleFunc(path, handler)

	return nil
}

func (server *WebServer) serveFiles(path string, staticPath string) {
	server.router.PathPrefix(path).Handler(http.FileServer(http.Dir(staticPath)))
}

func (server *WebServer) ListenAndServe() <-chan error {
	errs := make(chan error, 1)

	go func() {
		errs <- http.ListenAndServe(server.config.Addr, server.router)
		close(errs)
	}()

	return errs
}
