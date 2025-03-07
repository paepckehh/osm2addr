package osm2addr

import (
	"os"
	"sync"
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
	Country  string `json:"-"`
	City     string `json:"city"`
	Street   string `json:"street"`
	Postcode string `json:"postcode"`
}

// PLACE ...
type PLACE struct {
	CC []COUNTRY
}

// COUNTRY ..
type COUNTRY struct {
	Country string `json:"country"`
	PO      []POSTCODE
}

// POSTCODE
type POSTCODE struct {
	Postcode string `json:"postcode"`
	CI       []CITY
}

// CITY
type CITY struct {
	City string `json:"city"`
	ST   []STREET
}

// STREET
type STREET struct {
	Street string `json:"street"`
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
