package osm2addr

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/dustin/go-humanize"
)

// validDE CharacterMap
var validDE = [21]string{"Ä", "Ö", "Ü", "ä", "ö", "ü", "ß", "é", "è", "á", "ó", "ë", "ê", "š", "ł", "ć", "ú", "Ž", "Š"}

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
	h := humanize.Comma(int64(in))
	for {
		if len(h) < 10 {
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

// writeJsonFile ...
func writeJsonFile(countrycode, filename string, inMap map[string]map[string]bool) {
	folder := filepath.Join("json", countrycode)
	_ = os.MkdirAll(folder, 0755)
	j, err := json.Marshal(inMap)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(folder, filename), j, 0644); err != nil {
		panic(err)
	}
}
