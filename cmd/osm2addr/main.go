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

	// setup
	worker := 1 // runtime.NumCPU()
	file := "../../data/latest-germany.osm.pbf"

	// report
	fmt.Printf("\nOSM:Startup       # %v", ts)
	fmt.Printf("\nOSM:Worker        # %v", worker)
	fmt.Printf("\nOSM:File          # %v", file)

	// open file
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// parse file
	if err := osm2addr.Parse(f, worker); err != nil {
		log.Fatal(err)
	}

	// finish
	fmt.Printf("\nTotal Time Taken: %v\n", time.Now().Sub(ts))
}
