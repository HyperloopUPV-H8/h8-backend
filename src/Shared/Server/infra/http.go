package infra

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const serverAddr = "127.0.0.1:4000"

type HTTPServer[R, S any] struct {
	router     *mux.Router
	PacketRecv chan R
	PacketSend chan S
}

func New[R, S any](tx chan S, rx chan R) HTTPServer[R, S] {
	return HTTPServer[R, S]{
		router:     mux.NewRouter(),
		PacketRecv: rx,
		PacketSend: tx,
	}
}

func (server HTTPServer[R, S]) ListenAndServe() {
	log.Fatalln(http.ListenAndServe(serverAddr, server.router))
}
