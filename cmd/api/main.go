package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

  port := ":8080"

  fmt.Printf("Server running in http://localhost%s\n", port)

  if err := http.ListenAndServe(port, nil); err != nil {
    log.Fatal(err)
  }
}
