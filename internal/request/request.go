package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	state       int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	rawRequest, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	requestLine, numBytes, err := parseRequestLine(string(rawRequest))
	if err != nil {
		return nil, err
	}
	return &Request{RequestLine: requestLine}, nil
}

func parseRequestLine(request string) (RequestLine, int, error) {
	requestLine := strings.Split(request, "\r\n")[0]
	if !strings.Contains(requestLine, "\r\n") {
		return RequestLine{}, 0, nil
	}

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
	}, 0, nil
}

func isAllLetters(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
