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

type HTTPServer[D, O, M any] struct {
	router      *mux.Router
	page        spaHandler
	PacketRecv  chan D
	OrderSend   chan O
	MessageRecv chan M
}

func New[D, O, M any](dataIn chan D, orderOut chan O, messageIn chan M) HTTPServer[D, O, M] {
	return HTTPServer[D, O, M]{
		router:      mux.NewRouter(),
		page:        NewPage(defaultStaticPath, defaultIndexPath),
		PacketRecv:  dataIn,
		OrderSend:   orderOut,
		MessageRecv: messageIn,
	}
}

func (server HTTPServer[D, O, M]) ListenAndServe() {
	log.Fatalln(http.ListenAndServe(serverAddr, server.router))
}
