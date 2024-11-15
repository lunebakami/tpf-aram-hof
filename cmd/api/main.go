package main

import (
	"log"
	"tpf-aram-hof/cmd/server"
)

func main() {
  server := server.NewServer()

  if err := server.ListenAndServe(); err != nil {
    log.Fatal(err)
  }
}
