package infra

import (
	"net/http"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

func NewPage(staticPath string, indexPath string) spaHandler {
	return spaHandler{
		staticPath: staticPath,
		indexPath:  indexPath,
	}
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func (server HTTPServer[D, O, M]) HandleSPA() {
	server.router.PathPrefix("/").Handler(server.page)
}
