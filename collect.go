package osm2addr

import (
	"fmt"
	"strings"

	"github.com/agnivade/levenshtein"
)

// global
var seps = [...]string{" ", "-"}

// collect ...
func collect(target *Target) {

	// init
	var ok bool
	warning, warningCounter := make(map[string]int), 0
	corrected, correctedCounter := make(map[string]int), 0
	placeIDs := make(map[placeID]tagSet)
	places := make(map[postcode]map[city]map[street]placeIdHex)

	// range over targets channel
	for t := range targets {
		pid := id(string(t.Country) + string(t.Postcode) + string(t.City) + string(t.Street))
		if _, ok = placeIDs[pid]; ok {
			continue
		}
		if _, ok = places[t.Postcode]; !ok {
			places[t.Postcode] = make(map[city]map[street]placeIdHex)
			places[t.Postcode][t.City] = make(map[street]placeIdHex)
			switch t.Street {
			case "":
				continue
			default:
				places[t.Postcode][t.City] = make(map[street]placeIdHex)
				places[t.Postcode][t.City][t.Street] = pid.hex()
				placeIDs[pid] = *t
				e := fmt.Sprintf("[WARNING][NON-PRELOADED-POSTCODE-CITY-ADDED][POSTCODE:%v][CITY:%v]", t.Postcode, t.City)
				if _, ok = warning[e]; !ok {
					warning[e] = 0
				}
				warning[e]++
				warningCounter++
				continue
			}
		}
		if _, ok = places[t.Postcode][t.City]; !ok {
			for _, sep := range seps {
				if strings.Contains(string(t.City), sep) {
					c := strings.Split(string(t.City), sep)
					for ci := range places[t.Postcode] {
						if c[0] == string(ci) {
							e := fmt.Sprintf("[CORRECTED][CITY][SEP:%v]#%v#%v#", sep, t.City, ci)
							if _, ok = corrected[e]; !ok {
								corrected[e] = 0
							}
							corrected[e]++
							correctedCounter++
							t.City = ci
							pid = id(string(t.Country) + string(t.Postcode) + string(t.City) + string(t.Street))
							break
						}
					}
				}
			}
			places[t.Postcode][t.City] = make(map[street]placeIdHex)
			for ci := range places[t.Postcode] {
				if string(ci) == string(t.City) {
					continue
				}
				distance := levenshtein.ComputeDistance(string(ci), string(t.City))
				if distance < 2 {
					e := fmt.Sprintf("[WARNING][LEVENSHTEIN:%v][POSTCODE:%v][CITY]#%v#%v#", distance, t.Postcode, ci, t.City)
					if _, ok = warning[e]; !ok {
						warning[e] = 0
					}
					warning[e]++
					warningCounter++
				}
			}

		}
		switch t.Street {
		case "":
		default:
			if _, ok = places[t.Postcode][t.City][t.Street]; !ok {
				places[t.Postcode][t.City][t.Street] = pid.hex()
				placeIDs[pid] = *t
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
	writeJsonFile(target.Country, "placeID2addr.json", p)
	writeJsonFile(target.Country, "addr2placeID.json", places)
	writeJsonFile(target.Country, "warning.json", warning)
	writeJsonFile(target.Country, "corrected.json", corrected)
	collector.Done()
}
