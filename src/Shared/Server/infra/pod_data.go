package infra

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (server HTTPServer[D, O, M]) HandlePodData(route string, podData any) {
	server.router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handle pod data")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		r.Body.Close()

		encodedPodData, err := json.Marshal(podData)
		if err != nil {
			log.Fatalln(err)
		}

		w.Write(encodedPodData)
	})
}
