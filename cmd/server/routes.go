package server

import (
	"fmt"
	"net/http"
	"tpf-aram-hof/cmd/web"
	"tpf-aram-hof/cmd/web/hof"

	"github.com/a-h/templ"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("/assets/", fileServer)

	mux.Handle("/", templ.Handler(hof.HofBase()))
	mux.HandleFunc("/hof/player", hof.HofPostHandler)
  mux.HandleFunc("/hof/player/delete", hof.HofDeleteHandler)
	mux.HandleFunc("/hof/players", hof.HofGetHandler)

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	return mux
}

