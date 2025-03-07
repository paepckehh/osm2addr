package osm2addr

import (
	"fmt"
)

// place ...
type place struct {
	postcode string
	C        []city
}

// city ...
type city struct {
	name string
	S    []street
}

// street ...
type street struct {
	name string
}

// collect ...
func collect(targetCountry string) {

	// init
	var places []place
	placeIDs := make(map[ObjectID]bool)

	// range over targets channel
	for t := range targets {
		placeid := id(t.Country + t.Postcode + t.City + t.Street)
		if placeIDs[placeid] {
			continue // place found, skip
		}
		for _, place := range places {
			if place.postcode == t.Postcode {
				for _, city := range place.C {
					if city.name == t.City {
						for _, street := range city.S {
							if street.name == t.Street {
								panic("internal error, please report to developer team: placeIDs map-catcher faild") // unrechable
							}
						}
					}
				}
			}
		}
		placeIDs[placeid] = true
	}
	fmt.Printf("\nOSM:Collect:Places    # %v", hu(len(placeIDs)))
	collector.Done()
}

//				// levenshtein word distance
//				distance := levenshtein.ComputeDistance(city, t.City)
//				if distance < 2 {
//					fmt.Printf("\n[INFO] Levenshtein:%v \t# Postcode:%v \t# City:%v  \t<===> \tCity:%v", distance, t.Postcode, city, t.City)
//				}
//
//				// cologne phonetic collision
//				if phoneticCollision(city, t.City) {
//					fmt.Printf("\n[INFO] PhoneticCollision:1 \t# Postcode:%v \t# City:%v  \t<===> \tCity:%v", t.Postcode, city, t.City)
//				}

// write json mapping tables
//	writeJsonFile(targetCountry, "id.json", sets)
//	writeJsonFileMap(targetCountry, "city.json", c)
//	writeJsonFileMap(targetCountry, "street.json", s)
//	writeJsonFileMap(targetCountry, "postcode.json", p)
//	writeJsonFileMMap(targetCountry, "postcode2city.json", postcode2city)
//	writeJsonFileMMap(targetCountry, "postcode2city.json", postcode2city)
//	writeJsonFileMMap(targetCountry, "postcode2street.json", postcode2street)
//	writeJsonFileMMap(targetCountry, "city2postcode.json", city2postcode)
//	writeJsonFileMMap(targetCountry, "city2street.json", city2street)
//	fmt.Printf("\nOSM:Uniq:City         # %v", hu(len(c)))
//	fmt.Printf("\nOSM:Uniq:Street       # %v", hu(len(s)))
//	fmt.Printf("\nOSM:Uniq:Postcode     # %v", hu(len(p)))
//	fmt.Printf("\nOSM:Collect:Sets      # %v", hu(len(sets)))
