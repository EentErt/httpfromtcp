package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

const validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'*+-.^_`|~"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	dataString := string(data)
	if !strings.Contains(dataString, "\r\n") {
		return 0, false, nil
	}

	header := strings.Split(dataString, "\r\n")[0]
	if header == "" {
		return 0, true, nil
	}

	fields := strings.Fields(header)
	// check if there are two parts in the header
	if len(fields) != 2 {
		return 0, false, fmt.Errorf("invalid header format: %s", header)
	}

	// check if the header key ends with ":"
	if fields[0][len(fields[0])-1] != ':' {
		return 0, false, fmt.Errorf("invalid header format: %s", header)
	}
	for _, char := range fields[0][0 : len(fields[0])-1] {
		if !strings.ContainsRune(validChars, char) {
			return 0, false, fmt.Errorf("invalid character in header key: %s", fields[0][:len(fields[0])-1])
		}
	}
	key := strings.ToLower(fields[0][0 : len(fields[0])-1])

	value, ok := h[key]
	if ok {
		h[key] = value + ", " + fields[1]
	} else {
		h[key] = fields[1]
	}

	return len(header) + 2, false, nil
}

func NewHeaders() Headers {
	return make(Headers)
}
