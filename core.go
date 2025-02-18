package osm2addr

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"paepcke.de/osm2addr/internal/model"
	"paepcke.de/osm2addr/internal/pbf"
)

type AddrDB struct {
	CC []Country
}

type Country struct {
	CountryCode string
	DailCodes   []int32
	PO          []PostCodes
}

type PostCodes struct {
	PostCode string
	Cities   []string
	Streets  []string
}

func Parse(file *os.File, worker int) error {

	// init
	var wg sync.WaitGroup
	var globalErrCode error

	// osm decoder
	d, err := pbf.NewDecoder(context.Background(), file)
	if err != nil {
		return err
	}
	defer d.Close()

	// print osm file stats
	fmt.Printf("\nOSM:File:URL      # %v", d.Header.OsmosisReplicationBaseURL)
	fmt.Printf("\nOSM:File:Repl:USM # %v", d.Header.OsmosisReplicationSequenceNumber)
	fmt.Printf("\nOSM:File:Repl:TS  # %v\n", d.Header.OsmosisReplicationTimestamp)

	// spin up threads
	for i := 0; i < worker; i++ {
		wg.Add(1)
		go func() {
			w, tags, addrComplete, addrCompleteDE, nodes, objects := i, 0, 0, 0, 0, 0
			cc, country, countryErr := make(map[string]bool), 0, 0
			c, city, cityErr := make(map[string]bool), 0, 0
			s, street, streetErr := make(map[string]bool), 0, 0
			p, postcode, postcodeErr := make(map[string]bool), 0, 0
			postcode2city, postcode2street := make(map[string]map[string]bool), make(map[string]map[string]bool)
			city2postcode, city2street := make(map[string]map[string]bool), make(map[string]map[string]bool)
			for {
				objs, err := d.Decode()
				if err != nil {
					fmt.Printf("\naddr:country:uniq    %v  addr:country:valid %v  addr:country:err %v", hu(len(cc)), hu(country), hu(countryErr))
					fmt.Printf("\naddr:de:city:uniq    %v  addr:city:valid    %v  addr:city:err    %v", hu(len(c)), hu(city), hu(cityErr))
					fmt.Printf("\naddr:de:street:uniq  %v  addr:street:valid  %v  addr:street:err  %v", hu(len(s)), hu(street), hu(streetErr))
					fmt.Printf("\naddr:de:postcode:uniq%v  addr:postcode:valid%v  addr:postcode:err%v", hu(len(p)), hu(postcode), hu(postcodeErr))
					fmt.Printf("\naddr:de:complete     %v  addr:complete      %v", hu(addrCompleteDE), hu(addrComplete))
					fmt.Printf("\n\nWorker#%v Processed => Objects:%v => Nodes:%v => AddrTags:%v => ExitCode:%v\n", w, hu(objects), hu(nodes), hu(tags), err.Error())
					if err.Error() != "EOF" {
						globalErrCode = err
					}
					break
				}
				for _, obj := range objs {
					objects++
					switch obj.(type) {
					case *model.Node:
						nodes++
						n := obj.(*model.Node)
						if len(n.Tags) > 0 {
							tagCountry, tagPostcode, tagStreet, tagCity := "", "", "", ""
							// fmt.Printf("\n%v\n", n.Tags)
							for tag, content := range n.Tags {
								tags++
								if len(tag) > 8 {
									if tag[:5] == "addr:" {
										addrTag := strings.Split(tag, ":")
										if len(addrTag) > 1 {
											switch addrTag[1] {
											case "country":
												country++
												if len(content) != 2 {
													countryErr++
													continue
												}
												tagCountry = content
											case "street":
												street++
												l := len(content)
												if l > 3 && l > 256 {
													streetErr++
													continue
												}
												tagStreet = content
											case "city":
												city++
												l := len(content)
												if l > 3 && l > 256 {
													cityErr++
													continue
												}
												tagCity = content
											case "postcode":
												postcode++
												l := len(content)
												pInt, err := strconv.Atoi(content)
												if err != nil {
													postcodeErr++
													continue
												}
												pStr := strconv.Itoa(pInt)
												l = len(pStr)
												switch l {
												case 5:
												case 4:
													pStr = "0" + pStr
												default:
													postcodeErr++
													continue
												}
												tagPostcode = pStr
											}

										}
									}
								}
							}
							if tagCountry != "" {
								if !cc[tagCountry] {
									cc[tagCountry] = true
								}
								if tagPostcode != "" && tagStreet != "" && tagCity != "" {
									addrComplete++
									switch tagCountry {
									case "DE":
										// scope DE
										addrCompleteDE++
										if strings.Contains(tagCity, ".") {
											tagCity = tryNormaliseGermanCity(tagCity)
										}
										if !p[tagPostcode] {
											p[tagPostcode] = true
										}
										if !c[tagCity] {
											c[tagCity] = true
										}
										if !s[tagStreet] {
											s[tagStreet] = true
										}
										// scope DE:postcode
										// scope DE:postcode:city
										if _, ok := postcode2city[tagPostcode]; !ok {
											postcode2city[tagPostcode] = make(map[string]bool)
										}
										if !postcode2city[tagPostcode][tagCity] {
											postcode2city[tagPostcode][tagCity] = true
										}
										// scope DE:postcode:street
										if _, ok := postcode2street[tagPostcode]; !ok {
											postcode2street[tagPostcode] = make(map[string]bool)
										}
										if !postcode2street[tagPostcode][tagStreet] {
											postcode2street[tagPostcode][tagStreet] = true
										}
										// scope DE:city
										// scope DE:city:postcode
										if _, ok := city2postcode[tagCity]; !ok {
											city2postcode[tagCity] = make(map[string]bool)
										}
										if !city2postcode[tagCity][tagPostcode] {
											city2postcode[tagCity][tagPostcode] = true
										}
										// scope DE:city:street
										if _, ok := city2street[tagCity]; !ok {
											city2street[tagCity] = make(map[string]bool)
										}
										if !city2street[tagCity][tagStreet] {
											city2street[tagCity][tagStreet] = true
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
			writeJsonFile("postcode2city.json", postcode2city)
			writeJsonFile("postcode2street.json", postcode2street)
			writeJsonFile("city2postcode.json", city2postcode)
			writeJsonFile("city2street.json", city2street)
			wg.Done()
		}()
	}
	wg.Wait()
	return globalErrCode
}
