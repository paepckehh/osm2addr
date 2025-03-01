package osm2addr

import (
	"strings"
)

// camelCaseCityDE...
func camelCaseCityDE(in string) string {
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
	return out
}

// tryNormStreetDE ...
func tryNormStreetDE(in string) string {
	out := strings.ReplaceAll(in, "Str.", "Strasse")
	out = strings.ReplaceAll(out, "Straße", "Strasse")
	out = strings.ReplaceAll(out, "str.", "strasse")
	out = strings.ReplaceAll(out, "straße", "strasse")
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
	return out
}
