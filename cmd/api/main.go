package main

import (
	"log"
	"net/http"
	"time"

	"example.com/todo-ddd/internal/bootstrap"
)

func main() {
	module := bootstrap.Init()

	server := &http.Server{
		Addr:              ":8080",
		Handler:           module.HTTPHandler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("HTTP server listening on :8080")

	log.Fatal(server.ListenAndServe())
}
