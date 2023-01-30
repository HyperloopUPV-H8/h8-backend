package server

import (
	"log"
	"net/http"

	"github.com/HyperloopUPV-H8/Backend-H8/server/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Server struct {
	Router *mux.Router
}

func (server *Server) ServeData(route string, data []byte) {
	server.Router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		r.Body.Close()
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
}

func (server *Server) FileServer(route string, path string) {
	server.Router.PathPrefix(route).HandlerFunc(http.FileServer(http.Dir(path)).ServeHTTP)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (server *Server) WebsocketHandler(route string, handler models.WebsocketHandler) {
	server.Router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("server: WebsocketHandler: %s\n", err)
		}
		handler.HandleConn(conn)
	})
	log.Println(route)
}

func (server *Server) HandleFunc(route string, handler func(http.ResponseWriter, *http.Request)) {
	server.Router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		handler(w, r)
	})
}

func (server *Server) ListenAndServe(addr string) {
	http.ListenAndServe(addr, server.Router)
}
