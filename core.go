package osm2addr

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"paepcke.de/osm2addr/internal/model"
	"paepcke.de/osm2addr/internal/pbf"
)

// Target ...
type Target struct {
	Worker   int
	Country  string
	File     *os.File
	FileName string
}

// tagSET is a collection of osm tags:add<r
type tagSET struct {
	country  string
	city     string
	street   string
	postcode string
}

// Parse a Target
func Parse(target *Target) error {

	// init
	var wg sync.WaitGroup
	var globalErrCode error

	// osm decoder
	d, err := pbf.NewDecoder(context.Background(), target.File)
	if err != nil {
		return err
	}
	defer d.Close()

	// print osm file stats
	fmt.Printf("\nOSM:File:URL      # %v", d.Header.OsmosisReplicationBaseURL)
	fmt.Printf("\nOSM:File:Repl:USM # %v", d.Header.OsmosisReplicationSequenceNumber)
	fmt.Printf("\nOSM:File:Repl:TS  # %v\n", d.Header.OsmosisReplicationTimestamp)

	// spin up threads
	for i := 0; i < target.Worker; i++ {
		wg.Add(1)
		go func() {
			w, tags, addrComplete, addrCompleteCC, nodes, objects, uniformErr := i, 0, 0, 0, 0, 0, 0
			cc, country, countryErr := make(map[string]bool), 0, 0
			c, city, cityErr := make(map[string]bool), 0, 0
			s, street, streetErr := make(map[string]bool), 0, 0
			p, postcode, postcodeErr := make(map[string]bool), 0, 0
			postcode2city, postcode2street := make(map[string]map[string]bool), make(map[string]map[string]bool)
			city2postcode, city2street := make(map[string]map[string]bool), make(map[string]map[string]bool)
			for {
				objs, err := d.Decode()
				if err != nil {
					fmt.Printf("\naddr:country:uniq        %v  addr:all:country:valid %v  addr:country:err %v", hu(len(cc)), hu(country), hu(countryErr))
					fmt.Printf("\naddr:target:city:uniq    %v  addr:all:city:valid    %v  addr:city:err    %v", hu(len(c)), hu(city), hu(cityErr))
					fmt.Printf("\naddr:target:street:uniq  %v  addr:all:street:valid  %v  addr:street:err  %v", hu(len(s)), hu(street), hu(streetErr))
					fmt.Printf("\naddr:target:postcode:uniq%v  addr:all:postcode:valid%v  addr:postcode:err%v", hu(len(p)), hu(postcode), hu(postcodeErr))
					fmt.Printf("\naddr:target:records      %v  addr:all:records       %v  addr:uniform:err %v", hu(addrCompleteCC), hu(addrComplete), hu(uniformErr))
					fmt.Printf("\n\nWorker#%v Processed => Objects:%v => Nodes:%v => AddrTags:%v => ExitCode:%v\n", w, hu(objects), hu(nodes), hu(tags), err.Error())
					if err.Error() != "EOF" {
						globalErrCode = err
					}
					break
				}
				for _, obj := range objs {
					objects++
					switch o := obj.(type) {
					case *model.Node:
						nodes++
						if len(o.Tags) > 0 {
							t := tagSET{} // new tag set
							for tag, content := range o.Tags {
								tags++
								if len(tag) > 8 && tag[:5] == "addr:" {
									addrTag := strings.Split(tag, ":")
									if len(addrTag) > 1 {
										switch addrTag[1] {
										case "country":
											country++
											if len(content) != 2 {
												countryErr++
												continue
											}
											if content != strings.ToUpper(content) {
												countryErr++
												continue
											}
											t.country = content
										case "street":
											street++
											l := len(content)
											if l < 3 || l > 256 {
												streetErr++
												continue
											}
											t.street = content
										case "city":
											city++
											l := len(content)
											if l < 3 || l > 256 {
												cityErr++
												continue
											}
											t.city = content
										case "postcode":
											postcode++
											l := len(content)
											if l < 3 || l > 256 {
												postcodeErr++
												continue
											}
											t.postcode = content
										}

									}
								}
							}
							// validate tag complete, validate
							if t.country != "" {
								if !cc[t.country] {
									cc[t.country] = true // new country code found
								}
								if t.postcode != "" && t.street != "" && t.city != "" {
									addrComplete++
									if t.country == target.Country {
										addrCompleteCC++
										if err := t.uniform(); err != nil {
											uniformErr++
											continue
										}
										if !p[t.postcode] {
											p[t.postcode] = true
										}
										if !c[t.city] {
											c[t.city] = true
										}
										if !s[t.street] {
											s[t.street] = true
										}
										// scope CC:postcode
										// scope CC:postcode:city
										if _, ok := postcode2city[t.postcode]; !ok {
											postcode2city[t.postcode] = make(map[string]bool)
										}
										if !postcode2city[t.postcode][t.city] {
											postcode2city[t.postcode][t.city] = true
										}
										// scope CC:postcode:street
										if _, ok := postcode2street[t.postcode]; !ok {
											postcode2street[t.postcode] = make(map[string]bool)
										}
										if !postcode2street[t.postcode][t.street] {
											postcode2street[t.postcode][t.street] = true
										}
										// scope CC:city
										// scope CC:city:postcode
										if _, ok := city2postcode[t.city]; !ok {
											city2postcode[t.city] = make(map[string]bool)
										}
										if !city2postcode[t.city][t.postcode] {
											city2postcode[t.city][t.postcode] = true
										}
										// scope CC:city:street
										if _, ok := city2street[t.city]; !ok {
											city2street[t.city] = make(map[string]bool)
										}
										if !city2street[t.city][t.street] {
											city2street[t.city][t.street] = true
										}
									}
								}
							}
						}
					case *model.Way:
					case *model.Relation:
					default:
						panic("internal error, unknown osm model type")
					}
				}
			}
			// write json mapping tables
			writeJsonFile(target.Country, "postcode2city.json", postcode2city)
			writeJsonFile(target.Country, "postcode2street.json", postcode2street)
			writeJsonFile(target.Country, "city2postcode.json", city2postcode)
			writeJsonFile(target.Country, "city2street.json", city2street)
			wg.Done()
		}()
	}
	wg.Wait()
	return globalErrCode
}
