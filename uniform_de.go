package osm2addr

import (
	"fmt"
	"strconv"
	"strings"
)

// uniformDE
func (t *tagSet) uniformDE() int {
	count := 0
	if !isLatin1(string(t.City)) {
		fmt.Printf("\n[City][Latin1][%v]%v", t.Country, t.City)
		return 1
	}
	if t.Street != "" && !isLatin1(string(t.Street)) {
		fmt.Printf("\n[Street][Latin1][%v]%v", t.Country, t.Street)
		return 1
	}
	p, err := strconv.Atoi(string(t.Postcode))
	if err != nil {
		// fmt.Printf("[Postcode][%v]%v", t.Country, t.Postcode)
		return 1
	}
	pc := strconv.Itoa(p)
	switch len(pc) {
	case 5:
		t.Postcode = postcode(pc)
	case 4:
		t.Postcode = postcode("0" + pc)
	default:
		// fmt.Printf("[Postcode][%v]%v", t.Country, t.Postcode)
		return 1
	}
	var ok bool
	if t.City, ok = tryNormCityDE(t.City); !ok {
		count++
	}
	if t.Street != "" {
		if t.Street, ok = tryNormStreetDE(t.Street); !ok {
			count++
		}
	}
	return count
}

// tryNormStreetDE ...
func tryNormStreetDE(in street) (street, bool) {
	s := string(in)
	s = strings.ReplaceAll(s, "Strasse", "Straße")
	s = strings.ReplaceAll(s, "Str.", "Straße")
	s = strings.ReplaceAll(s, " str.", " Straße")
	s = strings.ReplaceAll(s, " strasse", " Straße")
	s = strings.ReplaceAll(s, "str.", " straße")
	s = strings.ReplaceAll(s, "strasse", "straße")
	if street(s) != in {
		// fmt.Printf("\n[UNIFORM][CITY][DE] IN-City:%v ======> OUT-City:%v", in, s)
		return street(s), false
	}
	return street(s), true
}

// tryNormCityDE ...
func tryNormCityDE(in city) (city, bool) {
	s := string(in)
	if strings.Contains(s, ".") {
		s = tryNormCityDEShortcut(s)
	}
	s = camelCaseCityDE(s)
	s = tryNormCityDETypo(s)
	if s != string(in) {
		// fmt.Printf("\n[UNIFORM][CITY][DE] IN-City:%v ======> OUT-City:%v", in, s)
		return city(s), false
	}
	return city(s), true
}

// tryNormCityDETypo ...
func tryNormCityDETypo(s string) string {
	s = strings.ReplaceAll(s, " von der ", " vor der ")
	s = strings.ReplaceAll(s, " in Isartal", " im Isartal")
	s = strings.ReplaceAll(s, " in Allgäu", " im Allgäu")
	s = strings.ReplaceAll(s, " in Rottal", " im Rottal")
	s = strings.ReplaceAll(s, " in Chiemgau", " im Chiemgau")
	s = strings.ReplaceAll(s, " in Grabfeld", " im Grabfeld")
	s = strings.ReplaceAll(s, "/ Saale", "/Saale")
	return s
}

// tryNormCityDEShortcut ...
func tryNormCityDEShortcut(s string) string {
	s = strings.ReplaceAll(s, " a.d. ", " an der ")
	s = strings.ReplaceAll(s, " a.d.", " an der ")
	s = strings.ReplaceAll(s, " a. d. ", " an der ")
	s = strings.ReplaceAll(s, " a. d.", " an der ")
	s = strings.ReplaceAll(s, " i.d. ", " in der ")
	s = strings.ReplaceAll(s, " i.d.", " in der ")
	s = strings.ReplaceAll(s, " i. d. ", " in der ")
	s = strings.ReplaceAll(s, " i. d.", " in der ")
	s = strings.ReplaceAll(s, " v.d. ", " von der ")
	s = strings.ReplaceAll(s, " v.d.", " von der ")
	s = strings.ReplaceAll(s, " v. d. ", " von der ")
	s = strings.ReplaceAll(s, " v. d.", " von der ")
	s = strings.ReplaceAll(s, " i. ", " in ")
	s = strings.ReplaceAll(s, " i.", " in ")
	s = strings.ReplaceAll(s, " a. ", " am ")
	s = strings.ReplaceAll(s, " a.", " am ")
	s = strings.ReplaceAll(s, " b. ", " bei ")
	s = strings.ReplaceAll(s, " b.", " bei ")
	return s
}

// camelCaseCityDE...
func camelCaseCityDE(s string) string {
	var out string
	lower := strings.ToLower(s)
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
