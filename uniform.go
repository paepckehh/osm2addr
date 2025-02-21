package osm2addr

import (
	"fmt"
	"strconv"
	"strings"
)

// uniform a tagSet
func (t *tagSET) uniform() error {
	switch t.country {
	case "DE":
		if !isLatin1(t.city) {
			return fmt.Errorf("[FAIL][City][Latin1][%v]%v", t.country, t.city)
		}
		if !isLatin1(t.street) {
			return fmt.Errorf("[FAIL][Street][Latin1][%v]%v", t.country, t.street)
		}
		if strings.Contains(t.city, ".") {
			t.city = tryNormCityDE(t.city)
		}
		t.city = camelCaseCityDE(t.city)
		p, err := strconv.Atoi(t.postcode)
		if err != nil {
			return fmt.Errorf("[FAIL][Postcode][%v]%v", t.country, t.postcode)
		}
		pc := strconv.Itoa(p)
		switch len(pc) {
		case 5:
			t.postcode = pc
		case 4:
			t.postcode = "0" + pc
		default:
			return fmt.Errorf("[FAIL][Postcode][%v]%v", t.country, t.postcode)
		}
	}
	return nil
}

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
