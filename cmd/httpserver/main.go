package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"httpfromtcp/internal/request"
	"httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	target := strings.Split(req.RequestLine.RequestTarget, "/")
	path := target[len(target)-1]
	switch path {
	case "yourproblem":
		return &server.HandlerError{
			StatusCode: 400,
			Message:    "Your problem is not my problem",
		}
	case "myproblem":
		return &server.HandlerError{
			StatusCode: 500,
			Message:    "Woopsie, my bad",
		}
	default:
		w.Write([]byte("All good, frfr"))
		return nil
	}
}
