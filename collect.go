package osm2addr

import (
	"fmt"
	"strings"

	"github.com/agnivade/levenshtein"
)

// collect ...
func collect(target *Target) {

	// init
	var ok bool
	warning, warningCounter := make(map[string]int), 0
	corrected, correctedCounter := make(map[string]int), 0
	placeIDs := make(map[placeID]tagSet)
	places := make(map[postcode]map[city]map[street]placeIdHex)
	places2 := make(map[postcode]map[city]map[street]bool)

	// range over targets channel
	for t := range targets {
		pid := id(string(t.Country) + string(t.Postcode) + string(t.City) + string(t.Street))
		if _, ok = placeIDs[pid]; ok {
			dbg("Drop:[Dedup][PlaceID][%v]%v %v %v", t.Country, t.Postcode, t.City, t.Street)
			continue
		}

		// new postcode: init maps and seed city
		if _, ok = places[t.Postcode]; !ok {
			places[t.Postcode] = make(map[city]map[street]placeIdHex)
			places2[t.Postcode] = make(map[city]map[street]bool)
			places[t.Postcode][t.City] = make(map[street]placeIdHex)
			places2[t.Postcode][t.City] = make(map[street]bool)
			if t.Street != "" {
				places[t.Postcode][t.City][t.Street] = pid.hex()
				places2[t.Postcode][t.City][t.Street] = true
				placeIDs[pid] = *t
				dbg("Add:[Place][NewPostcode][%v]%v %v %v", t.Country, t.Postcode, t.City, t.Street)
			}
			if !t.Preloaded {
				e := fmt.Sprintf("[WARNING][NON-PRELOADED-POSTCODE-CITY-ADDED][POSTCODE:%v][CITY:%v]", t.Postcode, t.City)
				warning[e]++
				warningCounter++
			}
			continue
		}

		// existing postcode, new city: merge only when the incoming city
		// differs from an existing one solely by separator choice (' ' vs '-'),
		// e.g. "Bad Homburg" vs "Bad-Homburg". A multi-token city must never be
		// collapsed into an unrelated single-token prefix ("Berlin Mitte" must
		// not merge into "Berlin"), otherwise its streets are silently lost.
		if _, ok = places[t.Postcode][t.City]; !ok {
			correctedCity := false
			inNorm := strings.ReplaceAll(string(t.City), "-", " ")
			for ci := range places[t.Postcode] {
				if inNorm == strings.ReplaceAll(string(ci), "-", " ") {
					e := fmt.Sprintf("[CORRECTED][CITY][SEP]#%v#%v#", t.City, ci)
					corrected[e]++
					correctedCounter++
					dbg("Norm:[City][Sep]#%v# =====> #%v#", t.City, ci)
					t.City = ci
					pid = id(string(t.Country) + string(t.Postcode) + string(t.City) + string(t.Street))
					correctedCity = true
					break
				}
			}
			if !correctedCity {
				places[t.Postcode][t.City] = make(map[street]placeIdHex)
				places2[t.Postcode][t.City] = make(map[street]bool)
				for ci := range places[t.Postcode] {
					distance := levenshtein.ComputeDistance(string(ci), string(t.City))
					if distance < 2 {
						e := fmt.Sprintf("[WARNING][LEVENSHTEIN:%v][POSTCODE:%v][CITY]#%v#%v#", distance, t.Postcode, ci, t.City)
						warning[e]++
						warningCounter++
						dbg("Warn:[Levenshtein:%v][Postcode:%v][City]#%v#%v#", distance, t.Postcode, ci, t.City)
					}
				}
			}
		}

		// re-check dedup after city correction may have changed pid
		if _, ok = placeIDs[pid]; ok {
			dbg("Drop:[Dedup][PostCorrect][PlaceID][%v]%v %v %v", t.Country, t.Postcode, t.City, t.Street)
			continue
		}

		// add street
		if t.Street != "" {
			if _, ok = places[t.Postcode][t.City][t.Street]; !ok {
				places[t.Postcode][t.City][t.Street] = pid.hex()
				places2[t.Postcode][t.City][t.Street] = true
				placeIDs[pid] = *t
				dbg("Add:[Place][%v]%v %v %v", t.Country, t.Postcode, t.City, t.Street)
			}
		}
	}
	fmt.Printf("\n----------------------------------------------------------------------------------")
	fmt.Printf("\nOSM:Corrected:Auto:Cases  # %v", hu(len(corrected)))
	fmt.Printf("\nOSM:Corrected:Auto:Total  # %v", hu(correctedCounter))
	fmt.Printf("\nOSM:Corrected:Warn:Cases  # %v", hu(len(warning)))
	fmt.Printf("\nOSM:Corrected:Warn:Total  # %v", hu(warningCounter))
	fmt.Printf("\nOSM:Collect:Places:Total  # %v", hu(len(placeIDs)))
	fmt.Printf("\n----------------------------------------------------------------------------------")
	p := make(map[string]tagSet, len(placeIDs))
	for pid, tset := range placeIDs {
		p[string(pid.hex())] = tset

	}
	writeJsonFile(target.Country, "addr.json", places2)
	writeJsonFile(target.Country, "addr2placeID.json", places)
	writeJsonFile(target.Country, "placeID2addr.json", p)
	writeJsonFile(target.Country, "warning.json", warning)
	writeJsonFile(target.Country, "corrected.json", corrected)
	collector.Done()
}
