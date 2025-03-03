package osm2addr

import (
	"context"
	"fmt"
	"strings"

	"github.com/f1monkey/phonetic/soundex"
	"paepcke.de/osm2addr/internal/model"
	"paepcke.de/osm2addr/internal/pbf"
)

// pbfparser for osm pbf files
func pbfparser(target *Target) {

	// osm decoder
	d, err := pbf.NewDecoder(context.Background(), target.File)
	if err != nil {
		panic(err)
	}
	defer d.Close()

	// print osm file stats
	fmt.Printf("\nOSM:PBF:File:URL      #  %v", d.Header.OsmosisReplicationBaseURL)
	fmt.Printf("\nOSM:PBF:File:Repl:USM #  %v", d.Header.OsmosisReplicationSequenceNumber)
	fmt.Printf("\nOSM:PBF:File:Repl:TS  #  %v", d.Header.OsmosisReplicationTimestamp)

	// init parser stats
	countries := make(map[string]bool)
	nodes, objects, tags, country, postcode, city, street, countryErr, cityErr, streetErr, postcodeErr, uniformErr, addrComplete := 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0

	for {
		objs, err := d.Decode()
		if err != nil {
			if err.Error() != "EOF" {
				panic(err)
			}
			fmt.Printf("\nOSM:PBF:ObjectsParsed # %v", hu(objects))
			fmt.Printf("\nOSM:PBF:AddrTags      # %v", hu(addrComplete))
			fmt.Printf("\nOSM:PBF:Uniq:Country  # %v", hu(len(countries)))
			fmt.Printf("\nOSM:PBF:Err:Uniform   # %v", hu(uniformErr))
			fmt.Printf("\nOSM:PBF:Err:Country   # %v", hu(countryErr))
			fmt.Printf("\nOSM:PBF:Err:Postcode  # %v", hu(postcodeErr))
			fmt.Printf("\nOSM:PBF:Err:City      # %v", hu(cityErr))
			fmt.Printf("\nOSM:PBF:Err:Street    # %v", hu(streetErr))
			parser.Done()
			return
		}
		for _, obj := range objs {
			objects++
			switch o := obj.(type) {
			case *model.Node:
				nodes++
				if len(o.Tags) > 0 {
					t := TagSET{} // init new tag set
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
									t.Country = content
								case "street":
									street++
									l := len(content)
									if l < 3 || l > 256 {
										streetErr++
										continue
									}
									t.Street = content
								case "city":
									city++
									l := len(content)
									if l < 3 || l > 256 {
										cityErr++
										continue
									}
									t.City = content
								case "postcode":
									postcode++
									l := len(content)
									if l < 3 || l > 256 {
										postcodeErr++
										continue
									}
									t.Postcode = content
								}

							}
						}
					}

					if t.Country != "" {
						if !countries[t.Country] {
							countries[t.Country] = true
						}
						if t.Postcode != "" && t.City != "" && t.Street != "" {
							addrComplete++
							if t.Country == target.Country {
								uniformErr = uniformErr + t.uniform()
								e := soundex.NewEncoder()
								t.CityPhonetic = e.Encode(t.City)
								targets <- &t
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
}
