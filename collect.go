package osm2addr

import "fmt"

// collect ...
func collect(targetCountry string) {

	// init
	sets := make(map[string]*TagSET)
	c, s, p := make(map[string]bool), make(map[string]bool), make(map[string]bool)
	postcode2city, postcode2street := make(map[string]map[string]bool), make(map[string]map[string]bool)
	city2postcode, city2street := make(map[string]map[string]bool), make(map[string]map[string]bool)

	// range over targets channel
	for t := range targets {
		tid := id(t.Country + t.Postcode + t.City + t.Street)
		tidb64 := tid.hex()
		if _, ok := sets[tidb64]; ok {
			continue
		}
		sets[tidb64] = t
		if !p[t.Postcode] {
			p[t.Postcode] = true
		}
		if !c[t.City] {
			c[t.City] = true
		}
		if !s[t.Street] {
			s[t.Street] = true
		}
		// scope CC:postcode
		// scope CC:postcode:city
		if _, ok := postcode2city[t.Postcode]; !ok {
			postcode2city[t.Postcode] = make(map[string]bool)
		}
		if !postcode2city[t.Postcode][t.City] {
			postcode2city[t.Postcode][t.City] = true
		}
		// scope CC:postcode:street
		if _, ok := postcode2street[t.Postcode]; !ok {
			postcode2street[t.Postcode] = make(map[string]bool)
		}
		if !postcode2street[t.Postcode][t.Street] {
			postcode2street[t.Postcode][t.Street] = true
		}
		// scope CC:city
		// scope CC:city:postcode
		if _, ok := city2postcode[t.City]; !ok {
			city2postcode[t.City] = make(map[string]bool)
		}
		if !city2postcode[t.City][t.Postcode] {
			city2postcode[t.City][t.Postcode] = true
		}
		// scope CC:city:street
		if _, ok := city2street[t.City]; !ok {
			city2street[t.City] = make(map[string]bool)
		}
		if !city2street[t.City][t.Street] {
			city2street[t.City][t.Street] = true
		}
	}
	// write json mapping tables
	writeJsonFile(targetCountry, "id.json", sets)
	writeJsonFileMap(targetCountry, "city.json", c)
	writeJsonFileMap(targetCountry, "street.json", s)
	writeJsonFileMap(targetCountry, "postcode.json", p)
	writeJsonFileMMap(targetCountry, "postcode2city.json", postcode2city)
	writeJsonFileMMap(targetCountry, "postcode2city.json", postcode2city)
	writeJsonFileMMap(targetCountry, "postcode2street.json", postcode2street)
	writeJsonFileMMap(targetCountry, "city2postcode.json", city2postcode)
	writeJsonFileMMap(targetCountry, "city2street.json", city2street)
	fmt.Printf("\nOSM:Uniq:City         # %v", hu(len(c)))
	fmt.Printf("\nOSM:Uniq:Street       # %v", hu(len(s)))
	fmt.Printf("\nOSM:Uniq:Postcode     # %v", hu(len(p)))
	fmt.Printf("\nOSM:Collect:Sets      # %v", hu(len(sets)))
	collector.Done()
}
