package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
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

func handler(w *response.Writer, req *request.Request) {
	target := strings.Split(req.RequestLine.RequestTarget, "/")
	path := target[len(target)-1]

	headers := response.GetDefaultHeaders(0)
	headers["content-type"] = "text/html"

	switch path {
	case "yourproblem":
		body := bodyBytes(400)
		w.WriteStatusLine(response.StatusNotFound)
		headers["content-length"] = strconv.Itoa(len(body))

		w.WriteHeaders(headers)
		w.WriteBody(body)
	case "myproblem":
		body := bodyBytes(500)
		w.WriteStatusLine(response.StatusServerError)
		headers["content-length"] = strconv.Itoa(len(body))

		w.WriteHeaders(headers)
		w.WriteBody(body)
	default:
		body := bodyBytes(200)
		w.WriteStatusLine(response.StatusOK)
		headers["content-length"] = strconv.Itoa(len(body))

		w.WriteHeaders(headers)
		w.WriteBody(body)

	}
}

func bodyBytes(code int) []byte {
	switch code {
	case 400:
		return []byte("<html>\n  <head>\n    <title>400 Bad Request</title>\n  </head>\n  <body>\n    <h1>Bad Request</h1>\n    <p>Your request honestly kinda sucked.</p>\n  </body>\n</html>")
	case 500:
		return []byte("<html>\n  <head>\n    <title>500 Internal Server Error</title>\n  </head>\n  <body>\n    <h1>Internal Server Error</h1>\n    <p>Okay, you know what? This one is on me.</p>\n  </body>\n</html>")
	default:
		return []byte("<html>\n  <head>\n    <title>200 OK</title>\n  </head>\n  <body>\n    <h1>Success!</h1>\n    <p>Your request was an absolute banger.</p>\n  </body>\n</html>")
	}
}

/*
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
*/
