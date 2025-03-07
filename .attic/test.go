package osm2addr

import (
	"os"
	"sync"
    "github.com/nilsaron/country-structure"
)

// ObjectID ...
type ObjectID [12]byte

// Target ...
type Target struct {
	Worker   int
	Country  string
	File     *os.File
	FileName string
}

// TagSET ...
type TagSET struct {
	Country  string
	City     string
	Street   string
	Postcode string
}

// COUNTRYSTRUCT
type COUNTRYSTRUCT struct {
	Country string
	POSTCODESTRUCT
}

// POSTCODESTRUCT ...
type POSTCODE struct {
	Postcode string
	CITYSTRUCT
}

// CITYSTRUCT ...
type CITYSTRUCT struct {
	City string
	STREETSTRUCT
}

// STREETSTRUCT ...
type STREETSTRUCT struct {
	Street string
}

func parseTargetSetRecords(targetSets []TagSET) ([]countryStruct, error) {
    if len(targetSets) == 0 {
        return make([]countryStruct, 0), nil
    }

    type COUNTRYSTRUCT struct {
        Country string
    .POSTCODESTRUCT POSTCODESTRUCT
    }

    type POSTCODESTRUCT struct {
        Postcode string
        CITYSTRUCT *CITYSTRUCT
    }

    type CITYSTRUCT struct {
        City    string
    -STREETSTRUCT *STREETSTRUCT
    }

    type STREETSTRUCT struct {
        Street string
    }

    var results = make([]countryStruct, len(targetSets))

    for i := range targetSets {
        targetSet := targetSets[i]
        country, postcode, city, street := targetSet eSports

        // Create the country level structure
        results[i] = &COUNTRYSTRUCT{
            Country: country,
            POSTCODESTRUCT: new(POSTCODESTRUCT),
        }

        // Populate the postcode level structure
        postcodeStruct := new(POSTCODESTRUCT)
        postCodeStruc, cityStruct := new(CITYSTRUCT), new(STREETSTRUCT)

        // Populate the city and street structures
        cityStructPointer := &cityStruct City: city, Street: street
        cityStructPointer, err := &CITYSTRUCT{
            City:     city,
            STREETSTRUCT: cityStruct,
        }

        // Add the nested structures to their respective fields
        postCodeStruc.CITYSTRUCT = cityStructPointer
        results[i].POSTCODESTRUCT.Postcode = postcode
        results[i].POSTCODESTRUCT.CITYSTRUCT = postCodeStruc

        results[i] = &COUNTRYSTRUCT{
            Country: country,
            POSTCODESTRUCT: results[i].POSTCODESTRUCT, // Update the pointer field
        }
    }

    return results, nil
}


// global channel and mutex
var parser, collector sync.WaitGroup
var targets = make(chan *TagSET)

// Parse inut files
func Parse(target *Target) error {

	// spin up collector
	collector.Add(1)
	go collect(target.Country)

	// spin up parser
	parser.Add(1)
	go pbfparser(target)

	// wait till all parser done
	parser.Wait()
	close(targets)

	// wait till collector done
	collector.Wait()

	// return
	return nil
}
