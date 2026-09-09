package osm2addr

import (
	"os"
	"sync"
)

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
		PostcodeLength int
	}
	PreLoadTrusted struct {
		File     *os.File
		Filename string
		Comma    rune
		Fields   int
		City     int
		Postcode int
		Street   int
	}

	// outDir overrides the json/ output base directory (used by tests)
	outDir string
}

// tagSet ...
type tagSet struct {
	Country   country  `json:"country"`
	Postcode  postcode `json:"postcode"`
	City      city     `json:"city"`
	Street    street   `json:"street"`
	Preloaded bool     `json:"-"`
}

type placeID [12]byte
type placeIdHex string
type country string
type postcode string
type city string
type street string

// outBase returns the directory that receives the per-country json output
func (t *Target) outBase() string {
	if t.outDir != "" {
		return t.outDir
	}
	return "json"
}

// Parse input files
func Parse(target *Target) error {

	// all producers feed the same unbuffered targets channel, the sole
	// consumer is the collector
	targets := make(chan *tagSet)
	var parser, collector sync.WaitGroup

	// spin up collector
	collector.Go(func() { collect(target, targets) })

	// trusted preload.csv feeds first, then the validated preload csv;
	// both run to completion before the pbf parser starts
	target.trustedPreloadFeed(targets)
	target.preloadFeed(targets)

	// spin up parser
	parser.Go(func() { pbfparser(target, targets) })

	// wait till parser done, then let the collector drain and exit
	parser.Wait()
	close(targets)
	collector.Wait()

	return nil
}
