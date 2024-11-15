package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"tpf-aram-hof/cmd/database"
)

type Server struct {
  port int

  db database.Service
}

func NewServer() *http.Server {
  port, _ := strconv.Atoi(os.Getenv("PORT"))
  NewServer := &Server {
    port: port,

    db: database.New(),
  }

  server := &http.Server{
    Addr: fmt.Sprintf(":%d", NewServer.port),
    Handler: NewServer.RegisterRoutes(),
    IdleTimeout: time.Minute,
    ReadTimeout: 10 * time.Second,
    WriteTimeout: 30 * time.Second,
  }

  fmt.Printf("Server is running on port %d\n", NewServer.port)

  return server
}
