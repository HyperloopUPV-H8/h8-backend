package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

type ServerConfig struct {
	Address            string
	FileServerPath     string `toml:"file_server_path"`
	FileServerEndpoint string `toml:"file_server_endpoint"`
	Endpoints          struct {
		FileServer       string `toml:"file_server"`
		PodData          string `toml:"pod_data"`
		OrderData        string `toml:"order_data"`
		UploadableBoards string `toml:"uploadable_boards"`
		Websocket        string
	}
}

type Server struct {
	router *mux.Router
	trace  zerolog.Logger
}

func New(router *mux.Router) *Server {
	trace.Info().Msg("new http server")
	return &Server{
		router: router,
		trace:  trace.With().Str("component", "httpServer").Logger(),
	}
}

func (server *Server) ServeData(route string, data any) {
	server.trace.Debug().Str("route", route).Msg("serve data")
	server.router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		r.Body.Close()
		w.Header().Set("Access-Control-Allow-Origin", "*")
		marshaledData, err := json.Marshal(data)
		if err != nil {
			server.trace.Error().Stack().Err(err).Msg("")
			http.Error(w, "failed to serialize resource", http.StatusInternalServerError)
			return
		}

		w.Write(marshaledData)
		server.trace.Trace().Str("route", route).Msg("write data")
	})
}

func (server *Server) FileServer(route string, path string) {
	server.trace.Debug().Str("route", route).Str("path", path).Msg("file server")
	server.router.PathPrefix(route).HandlerFunc(http.FileServer(http.Dir(path)).ServeHTTP)
}

func (server *Server) HandleFunc(route string, handler func(http.ResponseWriter, *http.Request)) {
	server.trace.Debug().Str("route", route).Msg("handle func")
	server.router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		server.trace.Trace().Str("route", route).Msg("handle request")
		handler(w, r)
	})
}

func (server *Server) ListenAndServe(addr string) {
	server.trace.Info().Str("addr", addr).Msg("listen and serve")
	http.ListenAndServe(addr, server.router)
}
