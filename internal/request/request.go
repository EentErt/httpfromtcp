package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

const bufferSize = 8

type state int

const (
	initialized state = iota
	done
)

type Request struct {
	RequestLine RequestLine
	state       state
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	buffer := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	request := Request{
		state: 0,
	}

	for request.state != done {
		if readToIndex == len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		bytesRead, err := reader.Read(buffer[readToIndex:])
		if err == io.EOF {
			request.state = done
			break
		} else if err != nil && err != io.EOF {
			return nil, fmt.Errorf("error reading from reader: %v", err)
		}

		readToIndex += bytesRead

		bytesParsed, err := request.parse(buffer[:readToIndex])
		if err != nil {
			return nil, fmt.Errorf("error parsing request: %v", err)
		} else if bytesParsed == 0 {
			continue
		}

		copy(buffer, buffer[bytesParsed:readToIndex])
		readToIndex -= bytesParsed
	}
	return &request, nil
}

func parseRequestLine(request string) (RequestLine, int, error) {
	if !strings.Contains(request, "\r\n") {
		return RequestLine{}, 0, nil
	}

	requestLine := strings.Split(request, "\r\n")[0]

	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return RequestLine{}, 0, fmt.Errorf("invalid request line: %s", requestLine)
	}

	httpVersion := strings.Split(parts[2], "/")[1]
	if httpVersion != "1.1" {
		return RequestLine{}, 0, fmt.Errorf("invalid HTTP version: %s", httpVersion)
	}

	// Check if all characters in the method are letters
	methodValid := isAllLetters(parts[0])

	//
	if strings.ToUpper(parts[0]) != parts[0] || !methodValid {
		return RequestLine{}, 0, fmt.Errorf("invalid method: %s", parts[0])
	}

	return RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: parts[1],
		Method:        parts[0],
	}, len(requestLine) + 2, nil
}

func isAllLetters(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case initialized:
		requestLine, numBytes, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if numBytes == 0 {
			return 0, nil
		}

		r.state = done
		r.RequestLine = requestLine
		return numBytes, nil
	case done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}
