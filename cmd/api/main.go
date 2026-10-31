package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"example.com/todo-ddd/internal/bootstrap"
)

func main() {
	start := time.Now()

	module := bootstrap.Init()

	server := &http.Server{
		Handler:           module.HTTPHandler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	startupMs := float64(time.Since(start).Microseconds()) / 1000

	log.Printf("Application started in %.2f ms", startupMs)
	log.Println("HTTP server listening on :8080")

	log.Fatal(server.Serve(listener))
}
