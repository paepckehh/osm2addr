package osm2addr

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const _secretMAC = "nachtssindallekatzenblau"

// convert a string into a repoduceable objectID
func id(in string) ObjectID {
	var o ObjectID
	h := sha256.Sum256([]byte(in + _secretMAC))
	copy(o[:], h[:12])
	return o
}

// convert a objectID into a hex string
func (in *ObjectID) hex() string {
	return hex.EncodeToString(in[:])
}

// isLatin1 ...
func isLatin1(s string) bool {
	if isASCII(s) {
		return true
	}
	var c = map[rune]rune{}
	for _, runeValue := range s {
		c[runeValue] = runeValue
	}
	for _, r := range c {
		if r > unicode.MaxLatin1 {
			return false
		}
	}
	return true
}

// isASCII ...
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// hu print large number readable for humans and fixed lenght
func hu(in int) string {
	p := message.NewPrinter(language.German)
	h := p.Sprintf("%d", in)
	for {
		if len(h) < 11 {
			h = " " + h
			continue
		}
		return h
	}
}

// makeCapitalLetter ....
func makeCapitalLetter(in string) string {
	if len(in) < 1 {
		return ""
	}
	r, size := utf8.DecodeRuneInString(in)
	if r == utf8.RuneError {
		panic("internal utf8 error")
	}
	return string(unicode.ToUpper(r)) + in[size:]
}

// camelCaseSeps ...
func camelCaseSeps(in string) string {
	out := camelCaseSep(in, "/")
	out = camelCaseSep(out, "-")
	out = camelCaseSep(out, ".")
	out = camelCaseSep(out, "(")
	return out
}

// camelCaseSep
func camelCaseSep(in, sep string) string {
	if strings.Contains(in, sep) {
		var out string
		parts := strings.Split(in, sep)
		for n, p := range parts {
			if n == 0 {
				out = p
				continue
			}
			p = makeCapitalLetter(p)
			out = out + sep + p
		}
		return out
	}
	return in
}

// containsSEP ..
func containsSEP(s string) bool {
	if strings.Contains(s, " ") || strings.Contains(s, "-") || strings.Contains(s, "/") {
		return true
	}
	return false
}
