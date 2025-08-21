package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"httpfromtcp/internal/headers"
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

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		proxyHandler(w, req)
		return
	}

	path := strings.TrimPrefix(req.RequestLine.RequestTarget, "/")

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

func proxyHandler(w *response.Writer, req *request.Request) {
	h := headers.Headers{
		"Content-Type":      "text/plain",
		"Transfer-Encoding": "chunked",
		"Trailer":           "X-Content-SHA256, X-Content-Length",
	}

	buffer := make([]byte, 32)
	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")

	resp, err := http.Get("https://httpbin.org" + target)
	if err != nil {
		w.WriteError(err)
	}
	defer resp.Body.Close()

	/*
		_, err = resp.Body.Read(buffer)
		if err != nil {
			w.WriteError(err)
		}
		defer resp.Body.Close()
	*/

	w.WriteStatusLine(response.StatusOK)

	w.WriteHeaders(h)

	body := []byte{}
	bodyLength := 0
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			body = append(body, buffer[:n]...)
			length, _ := w.WriteChunkedBody(buffer[:n])
			bodyLength += length
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			w.WriteError(err)
		}
	}
	w.WriteChunkedBodyDone()

	hash := sha256.Sum256(body)

	trailers := headers.Headers{
		"X-Content-SHA256": fmt.Sprintf("%x", hash),
		"X-Content-Length": strconv.Itoa(bodyLength),
	}

	w.WriteTrailers(trailers)
	w.WriteDone()
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
