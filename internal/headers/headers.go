package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Headers map[string]string

func NewHeaders() Headers {
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
	val := data[kidx+1 : crlIdx]

	if unicode.IsSpace(rune(key[len(key)-1])) {
		return 0, false, fmt.Errorf("bad header field name format: %s", string(data))
	}
	keyParts := bytes.Fields(key)
	if len( keyParts ) > 1 {
		return 0, false, fmt.Errorf("bad header field name format: %s", string(data))
	}

	keyString := strings.TrimSpace(string(key))
	valString := strings.TrimSpace(string(val))

	if !checkFieldName(keyString) {
		return 0, false, fmt.Errorf("bad header field name format: %s", string(keyString))
	}

	cleanKey := strings.ToLower(keyString)
	if _, ok := h[cleanKey]; ok {
		h[cleanKey] += fmt.Sprintf(", %s", valString)
	} else {
		h[cleanKey] = valString
	}

	return crlIdx + 2, false, nil
}

func (h Headers) Get(key string) (string, bool) {
	caser := cases.Title(language.English)
	val, ok := h[caser.String(key)]
	return val, ok
}
func (h Headers) Set(key string, values... string) {
	caser := cases.Title(language.English)
	cleanKey := caser.String(key)
	h[cleanKey] = values[0]
	if len(values) > 1 {
		for _, val := range values[1:] {
			h[cleanKey] = fmt.Sprintf("%s, %s", h[cleanKey], val) 
		}
	}
}
func (h Headers) Remove(key string) {
	caser := cases.Title(language.English)
	delete(h, caser.String(key))
}


func checkFieldName(fn string) bool {
	allowed := "!#$%&'*+-.^_`|~"

	for _, r := range fn {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && 
			!strings.ContainsAny(string(r), allowed) {
			return false
		}
	}

	return true
}
