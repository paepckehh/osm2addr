package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"paepcke.de/osm2addr"
)

var (
	version   = "dev"
	commit    = "none"
	buildtime = "unknown"
)

func main() {

	// init
	ts := time.Now()
	osm2addr.InitDebug()

	// setup defaults
	target := &osm2addr.Target{
		Worker:   1, // runtime.NumCPU()
		Country:  "DE",
		FileName: "data/germany-latest.osm.pbf",
	}

	// parse commandline options
	if len(os.Args) > 1 {
		if len(os.Args[1]) != 2 {
			log.Fatal("[OSM2ADDR][ERROR][FATAL] Target Country Code, if specified, must be two digits (example: osm2addr DE)")
		}
		target.Country = strings.ToUpper(os.Args[1])
	}
	if len(os.Args) > 2 {
		target.FileName = os.Args[2]
	}

	// report
	fmt.Printf("\nOSM:Version              # osm2addr %s (commit %s, built %s)", version, commit, buildtime)
	fmt.Printf("\nOSM:Startup               # %v", ts)
	fmt.Printf("\nOSM:TargetCountry         # %v", target.Country)
	fmt.Printf("\nOSM:WorkerScale           # %v", target.Worker)
	fmt.Printf("\nOSM:Debug                 # %v", osm2addr.DebugEnabled())
	fmt.Printf("\nOSM:File                  # %v", target.FileName)

	// open file
	t, err := os.Open(target.FileName)
	if err != nil {
		log.Fatalf("[OSM2ADDR][ERROR][FATAL] Unable to read: %v: %v", target.FileName, err)
	}
	target.File = t
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
