package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	StatusOK          StatusCode = 200
	StatusNotFound    StatusCode = 400
	StatusServerError StatusCode = 500
)

type writerState int

const (
	writeStatus writerState = iota
	writeHeaders
	writeBody
	writeTrailers
	writeDone
)

type Writer struct {
	Writer      io.Writer
	writerState writerState
}

func MakeWriter(writer io.Writer) *Writer {
	return &Writer{
		Writer:      writer,
		writerState: writeStatus,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writeStatus {
		return fmt.Errorf("cannot write to status line")
	}

	var err error

	switch statusCode {
	case StatusOK:
		_, err = w.Writer.Write(statusBytes(200, "OK"))
	case StatusNotFound:
		_, err = w.Writer.Write(statusBytes(400, "Not Found"))
	case StatusServerError:
		_, err = w.Writer.Write(statusBytes(500, "Server Error"))
	default:
		_, err = w.Writer.Write(statusBytes(500, ""))
	}

	if err != nil {
		return err
	}
	w.writerState = writeHeaders
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writeHeaders {
		return fmt.Errorf("cannot write headers")
	}
	for key, value := range headers {
		_, err := w.Writer.Write([]byte(key + ": " + value + "\r\n"))
		if err != nil {
			return err
		}
	}
	_, err := w.Writer.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	w.writerState = writeBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writeBody {
		return 0, fmt.Errorf("cannot write body")
	}
	i, err := w.Writer.Write(p)
	if err != nil {
		return 0, err
	}

	w.writerState = writeTrailers
	return i, nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != writeBody {
		return 0, fmt.Errorf("cannot write body")
	}

	lengthLine := fmt.Sprintf("%x\r\n", len(p))

	_, err := w.Writer.Write([]byte(lengthLine))
	if err != nil {
		return 0, err
	}

	j, err := w.Writer.Write(append(p, []byte("\r\n")...))
	if err != nil {
		return 0, err
	}

	return j - 2, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != writeBody {
		return 0, fmt.Errorf("cannot write body")
	}

	i, err := w.Writer.Write([]byte("0\r\n"))
	if err != nil {
		return 0, err
	}

	w.writerState = writeTrailers
	return i, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.writerState != writeTrailers {
		return fmt.Errorf("cannot write trailers")
	}
	for key, value := range h {
		_, err := w.Writer.Write([]byte(key + ": " + value + "\r\n"))
		if err != nil {
			return err
		}
	}
	_, err := w.Writer.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	w.writerState = writeDone
	return nil
}

func (w *Writer) WriteDone() error {
	_, err := w.Writer.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteError(err error) {
	headers := GetDefaultHeaders(0)
	body := []byte(fmt.Sprintf("%v", err))
	w.WriteStatusLine(StatusServerError)
	headers["content-length"] = strconv.Itoa(len(body))
	w.WriteHeaders(headers)
	w.WriteBody(body)
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		_, err := w.Write(statusBytes(200, "OK"))
		return err
	case StatusNotFound:
		_, err := w.Write(statusBytes(400, "Not Found"))
		return err
	case StatusServerError:
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
