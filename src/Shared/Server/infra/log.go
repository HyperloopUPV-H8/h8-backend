package infra

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func (server HTTPServer[D, O, M]) HandleLog(route string, loggerEnable chan<- bool) {
	server.router.Handle(route, LogHandle{logEnable: loggerEnable}).Methods(http.MethodPut)
}

type LogHandle struct {
	logEnable chan<- bool
}

func (handler LogHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle log")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("http server: log handle:", err)
	}

	go func() {
		if string(body) == "enable" {
			handler.logEnable <- true
		} else if string(body) == "disable" {
			handler.logEnable <- false
		}
	}()

	w.Write([]byte{})
}
