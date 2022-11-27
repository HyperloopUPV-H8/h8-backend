package infra

import (
	"log"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

var (
	serverAddr        = "127.0.0.1:4000"
	defaultIndexPath  = "index.html"
	defaultStaticPath = path.Join("static", "build")
)

type HTTPServer[D, O, M any] struct {
	router      *mux.Router
	page        spaHandler
	PacketChan  chan D
	OrderChan   chan O
	MessageChan chan M
}

func New[D, O, M any]() HTTPServer[D, O, M] {
	return HTTPServer[D, O, M]{
		router:      mux.NewRouter(),
		page:        NewPage(defaultStaticPath, defaultIndexPath),
		PacketChan:  make(chan D, 2024),
		OrderChan:   make(chan O, 2024),
		MessageChan: make(chan M, 2024),
	}
}

func (server HTTPServer[D, O, M]) ListenAndServe() {
	go log.Fatalln(http.ListenAndServe(serverAddr, server.router))
}
