package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"net"
	"strconv"
)

type Server struct {
	open     bool
	listener net.Listener
	handler  Handler
}

func Serve(port int, handlerFunc Handler) (*Server, error) {
	portString := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", portString)
	if err != nil {
		return nil, err
	}

	server := Server{
		open:     true,
		listener: listener,
		handler:  handlerFunc,
	}

	go server.listen()
	return &server, nil
}

func (s *Server) Close() error {
	if !s.open {
		return nil
	}

	s.open = false
	return s.listener.Close()
}

func (s *Server) listen() {
	for s.open {
		conn, err := s.listener.Accept()
		if err != nil {
			if !s.open {
				// server is not open, so we return
				return
			}
			// server is still open so we continue listening
			continue
		}

		s.handle(conn)
		conn.Close()
	}
}

func (s *Server) handle(conn net.Conn) {
	req, err := request.RequestFromReader(conn)
	if err != nil {
		writeError(conn, &HandlerError{StatusCode: 500, Message: fmt.Sprintf("Error: %v", err)})
	}

	buffer := bytes.Buffer{}

	handlerError := s.handler(&buffer, req)
	if handlerError.StatusCode != 200 {
		writeError(conn, handlerError)
	}

	if err := response.WriteStatusLine(conn, 200); err != nil {
		fmt.Println("Error writing response status line")
	}

	headers := response.GetDefaultHeaders(0)
	if err := response.WriteHeaders(conn, headers); err != nil {
		fmt.Printf("Error writing headers: %v\n", err)
	}

	// write the crlf between headers and body
	_, err = conn.Write([]byte("\r\n"))
	if err != nil {
		fmt.Printf("Error writing CRLF: %v\n", err)
	}

	// write the response body
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		fmt.Printf("Error writing CRLF: %v\n", err)
	}
}

type Handler func(w io.Writer, request *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func writeError(w io.Writer, h *HandlerError) {
	if err := response.WriteStatusLine(w, h.StatusCode); err != nil {
		fmt.Println("Error writing response status line")
	}

	writeBytes := []byte(string(int(h.StatusCode)) + "\r\n" + h.Message)
	w.Write(writeBytes)
}
