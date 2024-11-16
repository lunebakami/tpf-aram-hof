package server

import (
	"fmt"
	"net/http"
	"time"
	"tpf-aram-hof/cmd/web"
	"tpf-aram-hof/cmd/web/hello"
	"tpf-aram-hof/cmd/web/hof"

	"github.com/a-h/templ"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("/assets/", fileServer)
	mux.Handle("/web", templ.Handler(hello.HelloForm()))
	mux.HandleFunc("/hello", hello.HelloWebHandler)

  players := []hof.Player{
    {ID: 1, Nickname: "Player 1", Champion: "Champion 1", Description: "Description", GameMode: "ARAM",Frag: "Frag 1", Date: time.Now()},
    {ID: 2, Nickname: "Player 2", Champion: "Champion 2", Description: "Description", GameMode: "ARAM",Frag: "Frag 2", Date: time.Now()},
  }

	mux.Handle("/hof", templ.Handler(hof.HofBase()))
	mux.HandleFunc("/hof/player", hof.HofWebHandler)
	mux.Handle("/hof/players", templ.Handler(hof.HofList(players)))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	return mux
}
