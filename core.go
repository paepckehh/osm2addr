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
	PreLoad  struct {
		File           *os.File
		Filename       string
		Fields         int
		City           int
		Postcode       int
		PostcodeLenght int
	}
}

type country string  // `json:"postcode"`
type postcode string // `json:"postcode"`
type city string     // `json:"city"`
type street string   // `json:"street"`

// TagSET ...
type TagSET struct {
	Country  country  // `json:"-"`
	Postcode postcode // `json:"postcode"`
	City     city     // `json:"city"`
	Street   street   // `json:"street"`
}

// global channel and mutex
var preload, parser, collector sync.WaitGroup
var targets = make(chan *TagSET)

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
