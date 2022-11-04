package infra

import (
	"net/http"
	"os"
	"path/filepath"
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
	path := filepath.Join(h.staticPath, r.URL.Path)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
