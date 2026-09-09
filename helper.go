package osm2addr

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const _secretMAC = "nachtssindallekatzenblau"

// convert a string into a repoduceable objectID
func id(in string) placeID {
	var o placeID
	h := sha256.Sum256([]byte(in + _secretMAC))
	copy(o[:], h[:12])
	return o
}

// convert a objectID into a hex string
func (in *placeID) hex() placeIdHex {
	return placeIdHex(hex.EncodeToString(in[:]))
}

// isLatin1 ...
func isLatin1(s string) bool {
	if isASCII(s) {
		return true
	}
	for _, r := range s {
		if r > unicode.MaxLatin1 {
			return false
		}
	}
	return true
}

// isASCII ...
func isASCII(s string) bool {
	for i := range len(s) {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// hu print large number readable for humans and fixed lenght
func hu(in int) string {
	p := message.NewPrinter(language.German)
	return fmt.Sprintf("%11s", p.Sprintf("%d", in))
}

// makeCapitalLetter ....
func makeCapitalLetter(in string) string {
	r, size := utf8.DecodeRuneInString(in)
	if r == utf8.RuneError && size <= 1 {
		return in
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
