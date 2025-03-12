package osm2addr

import (
	"os"
	"sync"
)

// Schema ...
type placeID [12]byte
type placeIdHex string //`json:"placeid"`
type country string    //`json:"country"`
type postcode string   //`json:"postcode"`
type city string       //`json:"city"`
type street string     //`json:"street"`

// Target ...
type Target struct {
	Worker   int
	Country  string
	File     *os.File
	FileName string
	PreLoad  struct {
		File           *os.File
		Filename       string
		Fields         int
		City           int
		Postcode       int
		PostcodeLenght int
	}
}

// tagSet ...
type tagSet struct {
	Country  country  `json:"country"`
	Postcode postcode `json:"postcode"`
	City     city     `json:"city"`
	Street   street   `json:"street"`
}

// global channel and mutex
var preload, parser, collector sync.WaitGroup
var targets = make(chan *tagSet)

// Parse inut files
func Parse(target *Target) error {

	// checkPreloadFile
	preload.Add(1)
	go target.preloadFeed()

	// spin up collector
	collector.Add(1)
	go collect(target)

	// wait till preload is done
	preload.Wait()

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
