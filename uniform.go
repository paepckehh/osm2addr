package osm2addr

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// uniform tag
func (t *osmTag) uniform() {
	switch t.country {
	case "DE":
		if strings.Contains(t.city, ".") {
			t.city = tryNormCityDE(t.city)
		}
		t.city = camelCase(t.city)
		if len(t.postcode) == 4 {
			t.postcode = "0" + t.postcode
		}
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

// camelCase...
func camelCase(in string) string {
	var out string
	lower := strings.ToLower(in)
	parts := strings.Split(lower, " ")
	for n, p := range parts {
		if n == 0 {
			out = camelCaseSeps(makeCapitalLetter(p))
			continue
		}
		switch len(p) {
		case 0:
		case 1:
		default:
			switch p {
			case "an":
			case "am":
			case "auf":
			case "bei":
			case "beim":
			case "der":
			case "den":
			case "dem":
			case "ii":
				p = "II"
			case "in":
			case "im":
			case "kalten":
				p = "kalten"
			case "unter":
			case "ob":
			case "ot":
				p = "OT"
			case "ofr.":
				p = "OFr."
			case "opf.":
				p = "OPf."
			case "nb":
				p = "NB"
			case "von":
			case "vor":
			case "vorm":
			default:
				p = camelCaseSeps(makeCapitalLetter(p))
			}
		}
		out = out + " " + p
	}
	if out != in {
		fmt.Printf("\n[CamelCased]%v=>%v", in, out)
	}
	return out
}

// tryNormCityDE ...
func tryNormCityDE(in string) string {
	out := strings.ReplaceAll(in, " a.d. ", " an der ")
	out = strings.ReplaceAll(out, " a.d.", " an der ")
	out = strings.ReplaceAll(out, " a. d. ", " an der ")
	out = strings.ReplaceAll(out, " a. d.", " an der ")
	out = strings.ReplaceAll(out, " i.d. ", " in der ")
	out = strings.ReplaceAll(out, " i.d.", " in der ")
	out = strings.ReplaceAll(out, " i. d. ", " in der ")
	out = strings.ReplaceAll(out, " i. d.", " in der ")
	out = strings.ReplaceAll(out, " v.d. ", " von der ")
	out = strings.ReplaceAll(out, " v.d.", " von der ")
	out = strings.ReplaceAll(out, " v. d. ", " von der ")
	out = strings.ReplaceAll(out, " v. d.", " von der ")
	out = strings.ReplaceAll(out, " i. ", " in ")
	out = strings.ReplaceAll(out, " i.", " in ")
	out = strings.ReplaceAll(out, " a. ", " am ")
	out = strings.ReplaceAll(out, " a.", " am ")
	out = strings.ReplaceAll(out, " b. ", " bei ")
	out = strings.ReplaceAll(out, " b.", " bei ")
	// if out != in { fmt.Printf("\n[CITY-DE-FORMED]%v=>%v", in, out) }
	return out
}
