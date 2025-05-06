package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers{
	return make(map[string]string)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	kidx := bytes.Index(data, []byte(":"))
	crlIdx := bytes.Index(data, []byte("\r\n"))

	if crlIdx == -1 {
		return 0, false, nil
	}
	if crlIdx == 0 {
		return 2, true, nil
	}

	if kidx == -1 {
		return 0, false, fmt.Errorf("bad header format: %s", string(data))
	}

	key := data[:kidx]
	val := data[kidx+1:crlIdx]
	
	if unicode.IsSpace(rune(key[len(key)-1])) {
		return 0, false, fmt.Errorf("bad header key format: %s", string(data))
	}

	keyString := strings.TrimSpace(string(key))
	valString := strings.TrimSpace(string(val))

	h[keyString] = valString

	return crlIdx+2, false, nil
}
