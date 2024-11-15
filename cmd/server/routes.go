package server

import (
	"fmt"
	"net/http"
	"time"
)

type Player struct {
	ID       int       `json:"id"`
	Nickname string    `json:"nickname"`
	Champion string    `json:"champion"`
	Frag     string    `json:"frag"`
	Date     time.Time `json:"date"`
}

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	return mux
}
