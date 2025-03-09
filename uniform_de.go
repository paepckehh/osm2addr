package osm2addr

import (
	"strconv"
	"strings"
)

// uniformDE
func (t *TagSET) uniformDE() int {
	count := 0
	if !isLatin1(string(t.City)) {
		// fmt.Printf("[City][Latin1][%v]%v", t.Country, t.City)
		return 1
	}
	if !isLatin1(string(t.Street)) {
		// fmt.Printf("[Street][Latin1][%v]%v", t.Country, t.Street)
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
	switch p {
	case 99334:
		if strings.Contains(string(t.City), "Elxleben") {
			t.City = "Elxleben am Steiger"
			count++
		}
	case 25761:
		if strings.Contains(string(t.City), "Westerdeichstrich") {
			t.City = "Westerdeichstrich (Kreis Dithmarschen)"
			count++
		}
	case 25862:
		if strings.Contains(string(t.City), "Goldelund") {
			t.City = "Goldelund Nordfriesland"
			count++
		}
	case 93453:
		if strings.Contains(string(t.City), "Neukirchen") {
			t.City = "Neukirchen b.Hl.Blut"
			count++
		}
	}
	var ok bool
	if ok, t.City = tryNormCityDE(t.City); !ok {
		count++
	}
	if ok, t.Street = tryNormStreetDE(t.Street); !ok {
		count++
	}
	return count
}

// tryNormStreetDE ...
func tryNormStreetDE(in street) (bool, street) {
	s := string(in)
	s = strings.ReplaceAll(s, "Strasse", "Straße")
	s = strings.ReplaceAll(s, "Str.", "Straße")
	s = strings.ReplaceAll(s, " str.", " Straße")
	s = strings.ReplaceAll(s, " strasse", " Straße")
	s = strings.ReplaceAll(s, "str.", " straße")
	s = strings.ReplaceAll(s, "strasse", "straße")
	if street(s) != in {
		// fmt.Printf("\n[UNIFORM][CITY][DE] IN-City:%v ======> OUT-City:%v", in, s)
		return false, street(s)
	}
	return true, street(s)
}

// tryNormCityDE ...
func tryNormCityDE(in city) (bool, city) {
	s := string(in)
	if strings.Contains(s, ".") {
		s = tryNormCityDEShortcut(s)
	}
	s = camelCaseCityDE(s)
	s = tryNormCityDETypo(s)
	if s != string(in) {
		// fmt.Printf("\n[UNIFORM][CITY][DE] IN-City:%v ======> OUT-City:%v", in, s)
		return false, city(s)
	}
	return true, city(s)
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
