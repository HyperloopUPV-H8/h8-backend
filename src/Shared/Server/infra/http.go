package infra

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	serverAddr        = "127.0.0.1:4000"
	defaultIndexPath  = "index.html"
	defaultStaticPath = ""
)

type HTTPServer[R, S any] struct {
	router     *mux.Router
	page       spaHandler
	PacketRecv chan R
	OrderSend  chan S
}

func New[R, S any](rx chan R, tx chan S) HTTPServer[R, S] {
	return HTTPServer[R, S]{
		router:     mux.NewRouter(),
		page:       NewPage(defaultStaticPath, defaultIndexPath),
		PacketRecv: rx,
		OrderSend:  tx,
	}
}

func (server HTTPServer[R, S]) ListenAndServe() {
	log.Fatalln(http.ListenAndServe(serverAddr, server.router))
}
