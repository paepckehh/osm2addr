package osm2addr

import (
	"fmt"
	"strconv"
	"strings"
)

// uniform a tagSet
func (t *TagSET) uniform() error {
	switch t.Country {
	case "DE":
		if !isLatin1(t.City) {
			return fmt.Errorf("[FAIL][City][Latin1][%v]%v", t.Country, t.City)
		}
		if !isLatin1(t.Street) {
			return fmt.Errorf("[FAIL][Street][Latin1][%v]%v", t.Country, t.Street)
		}
		p, err := strconv.Atoi(t.Postcode)
		if err != nil {
			return fmt.Errorf("[FAIL][Postcode][%v]%v", t.Country, t.Postcode)
		}
		pc := strconv.Itoa(p)
		switch len(pc) {
		case 5:
			t.Postcode = pc
		case 4:
			t.Postcode = "0" + pc
		default:
			return fmt.Errorf("[FAIL][Postcode][%v]%v", t.Country, t.Postcode)
		}
		if strings.Contains(t.City, ".") {
			t.City = tryNormCityDE(t.City)
		}
		t.City = camelCaseCityDE(t.City)
		t.Street = tryNormStreetDE(t.Street)
	}
	return nil
}
