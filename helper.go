package osm2addr

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const _secretMAC = "Nachts99Sind44Alle66Grau99Blau"

// convert a string into a repoduceable objectID
func id(in string) ObjectID {
	var o ObjectID
	h := sha256.Sum256([]byte(in + _secretMAC))
	copy(o[:], h[:12])
	return o
}

// convert a objectID into a base64 string
func (in *ObjectID) b64() string {
	return base64.StdEncoding.EncodeToString(in[:])
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
func writeJsonFile(countrycode, filename string, in map[string]TagSET) {
	folder := filepath.Join("json", countrycode)
	_ = os.MkdirAll(folder, 0755)
	j, err := json.MarshalIndent(&in, "", "\t")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(folder, filename), j, 0644); err != nil {
		panic(err)
	}
}

// writeJsonFileIDMap ...
//func writeJsonFileIDMap(countrycode, filename string, in map[string]objectID) {
//	folder := filepath.Join("json", countrycode)
//	_ = os.MkdirAll(folder, 0755)
//	out := make(map[string]string, len(in))
//	for key, value := range in {
//		out[key] = value.b64()
//	}
//	j, err := json.MarshalIndent(&out, "", "\t")
//	if err != nil {
//		panic(err)
//	}
//	if err := os.WriteFile(filepath.Join(folder, filename), j, 0644); err != nil {
//		panic(err)
//	}
//}

// writeJsonFileMap ...
func writeJsonFileMap(countrycode, filename string, in map[string]bool) {
	folder := filepath.Join("json", countrycode)
	_ = os.MkdirAll(folder, 0755)
	j, err := json.MarshalIndent(&in, "", "\t")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(folder, filename), j, 0644); err != nil {
		panic(err)
	}
}

// writeJsonFileMMap ...
func writeJsonFileMMap(countrycode, filename string, in map[string]map[string]bool) {
	folder := filepath.Join("json", countrycode)
	_ = os.MkdirAll(folder, 0755)
	j, err := json.MarshalIndent(&in, "", "\t")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(folder, filename), j, 0644); err != nil {
		panic(err)
	}
}
