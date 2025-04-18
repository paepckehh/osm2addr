package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"paepcke.de/osm2addr"
)

func main() {

	// init
	ts := time.Now()
	var err error

	// setup defaults
	target := &osm2addr.Target{
		Worker:   1, // runtime.NumCPU()
		Country:  "DE",
		FileName: "../../data/germany-latest.osm.pbf",
	}

	// parse commandline options
	if len(os.Args) > 1 {
		if len(os.Args[1]) != 2 {
			log.Fatal("[OSM2ADDR][ERROR][FATAL] Target Country Code, if specified, must be two digits (example: osm2addr DE)")
		}
		target.Country = os.Args[1]
	}
	if len(os.Args) > 2 {
		t, err := os.Open(os.Args[2])
		if err != nil {
			log.Printf("[OSM2ADDR][ERROR][FATAL] Unable to read: %v", os.Args[2])
			log.Fatal("[OSM2ADDR][ERROR][FATAL] Target File, if specified, must be a readable file (example: osm2addr DE ../../data/germany-latest.osm.pbf)")
		}
		if err := t.Close(); err != nil {
			log.Printf("[OSM2ADDR][ERROR][FATAL] Unable to cloe: %v: %v", os.Args[2], err)
		}
		target.FileName = os.Args[2]
	}

	// report
	fmt.Printf("\nOSM:Startup               # %v", ts)
	fmt.Printf("\nOSM:TargetCountry         # %v", target.Country)
	fmt.Printf("\nOSM:WorkerScale           # %v", target.Worker)
	fmt.Printf("\nOSM:File                  # %v", target.FileName)

	// open file
	target.File, err = os.Open(target.FileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := target.File.Close(); err != nil {
			fmt.Printf("[OSM2ADDR][ERROR] during file close: %s", err)
		}
	}()

	// parse file
	if err := osm2addr.Parse(target); err != nil {
		log.Fatal(err)
	}

	// finish
	fmt.Printf("\nOSM:Time:Total            # %v\n", time.Since(ts))
}
