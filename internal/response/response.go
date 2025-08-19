package response

import (
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	statusOK          StatusCode = 200
	statusNotFound    StatusCode = 400
	statusServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case statusOK:
		_, err := w.Write(statusBytes(200, "OK"))
		return err
	case statusNotFound:
		_, err := w.Write(statusBytes(400, "Not Found"))
		return err
	case statusServerError:
		_, err := w.Write(statusBytes(500, "Server Error"))
		return err
	default:
		_, err := w.Write(statusBytes(500, ""))
		return err
	}
}

func statusBytes(statusCode int, reason string) []byte {
	codeString := strconv.Itoa(statusCode)
	return []byte("HTTP/1.1 " + codeString + " " + reason + "\r\n")
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	return headers.Headers{
		"content-length": strconv.Itoa(contentLen),
		"connection":     "close",
		"content-type":   "text/plain",
	}
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := w.Write([]byte(key + ": " + value + "\r\n"))
		if err != nil {
			return err
		}
	}

	return nil
}
