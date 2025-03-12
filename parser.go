package osm2addr

import (
	"context"
	"fmt"
	"strings"

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

	// print startup info
	fmt.Printf("\n----------------------------------------------------------------------------------")
	fmt.Printf("\nOSM:PBF:File:URL          # %v", d.Header.OsmosisReplicationBaseURL)
	fmt.Printf("\nOSM:PBF:File:Repl:USM     # %v", d.Header.OsmosisReplicationSequenceNumber)
	fmt.Printf("\nOSM:PBF:File:Repl:TS      # %v", d.Header.OsmosisReplicationTimestamp)

	// init parser stats
	countries := make(map[country]bool)
	objectCounter, tagCounter, countryCounter, postcodeCounter, cityCounter, streetCounter := 0, 0, 0, 0, 0, 0
	countryErr, cityErr, streetErr, postcodeErr, uniformErr, addrComplete := 0, 0, 0, 0, 0, 0

	for {
		objs, err := d.Decode()
		if err != nil {
			if err.Error() != "EOF" {
				panic(err)
			}
			fmt.Printf("\n----------------------------------------------------------------------------------")
			fmt.Printf("\nOSM:PBF:Parsed:Objects    # %v", hu(objectCounter))
			fmt.Printf("\nOSM:PBF:Parsed:Tags       # %v", hu(tagCounter))
			fmt.Printf("\nOSM:PBF:Parsed:Country    # %v", hu(countryCounter))
			fmt.Printf("\nOSM:PBF:Parsed:Street     # %v", hu(streetCounter))
			fmt.Printf("\nOSM:PBF:Parsed:City       # %v", hu(cityCounter))
			fmt.Printf("\nOSM:PBF:Parsed:Postcode   # %v", hu(postcodeCounter))
			fmt.Printf("\n----------------------------------------------------------------------------------")
			fmt.Printf("\nOSM:PBF:Complete:AddrTags # %v", hu(addrComplete))
			fmt.Printf("\nOSM:PBF:Uniq:Country      # %v", hu(len(countries)))
			fmt.Printf("\nOSM:PBF:Err:Uniform       # %v", hu(uniformErr))
			fmt.Printf("\nOSM:PBF:Err:Country       # %v", hu(countryErr))
			fmt.Printf("\nOSM:PBF:Err:Postcode      # %v", hu(postcodeErr))
			fmt.Printf("\nOSM:PBF:Err:City          # %v", hu(cityErr))
			fmt.Printf("\nOSM:PBF:Err:Street        # %v", hu(streetErr))
			parser.Done()
			return
		}
		for _, obj := range objs {
			objectCounter++
			switch o := obj.(type) {
			case *model.Node:
				if len(o.Tags) > 0 {
					t := tagSet{} // init new tag set
					for tag, content := range o.Tags {
						tagCounter++
						if len(tag) > 8 && tag[:5] == "addr:" {
							addrTag := strings.Split(tag, ":")
							if len(addrTag) > 1 {
								switch addrTag[1] {
								case "country":
									countryCounter++
									if len(content) != 2 {
										countryErr++
										continue
									}
									if content != strings.ToUpper(content) {
										countryErr++
										continue
									}
									t.Country = country(content)
								case "street":
									streetCounter++
									l := len(content)
									if l < 3 || l > 256 {
										streetErr++
										continue
									}
									t.Street = street(content)
								case "city":
									cityCounter++
									l := len(content)
									if l < 2 || l > 256 {
										cityErr++
										continue
									}
									t.City = city(content)
								case "postcode":
									postcodeCounter++
									l := len(content)
									if l < 3 || l > 256 {
										postcodeErr++
										continue
									}
									t.Postcode = postcode(content)
								}

							}
						}
					}

					if t.Country != "" {
						if !countries[country(t.Country)] {
							countries[country(t.Country)] = true
						}
						if t.Postcode != "" && t.City != "" && t.Street != "" {
							addrComplete++
							if t.Country == country(target.Country) {
								uniformErr += t.uniform()
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
