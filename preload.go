package osm2addr

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// preloadFeed ...
func (target *Target) preloadFeed() {
	if target.checkPreloadFile() {

		// init
		fail, failCounter, counter := make(map[string]int), 0, 0
		defer func() {
			if err := target.PreLoad.File.Close(); err != nil {
				fmt.Printf("[OSM2ADDR][ERROR] unable to close file: %s", err)
			}
		}()

		// report
		fmt.Printf("\n----------------------------------------------------------------------------------")
		fmt.Printf("\nOSM:PreLoadFile           # %v", target.PreLoad.Filename)

		// setup
		r := csv.NewReader(target.PreLoad.File)
		r.Comma = ','
		r.FieldsPerRecord = target.PreLoad.Fields
		r.LazyQuotes = false
		r.ReuseRecord = true
		_, _ = r.Read() // skip head

		// loop
		for {
			row, err := r.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				e := fmt.Sprintf("[ERROR][PRE][CSV][PARSE]#%v#", err)
				fail[e]++
				failCounter++
				dbg("Drop:[Preload][CSV][Parse]%v", err)
				continue
			}
			if len(row[target.PreLoad.City]) < 2 {
				e := fmt.Sprintf("[ERROR][PRE][CITY][LENGHT]#%v#%v#", row[target.PreLoad.Postcode], row[target.PreLoad.City])
				fail[e]++
				failCounter++
				dbg("Drop:[Preload][City][Length]%v %v", row[target.PreLoad.Postcode], row[target.PreLoad.City])
				continue
			}
			if len(row[target.PreLoad.Postcode]) != target.PreLoad.PostcodeLenght {
				e := fmt.Sprintf("[ERROR][PRE][POSTCODE][LENGHT]#%v#%v#", row[target.PreLoad.Postcode], row[target.PreLoad.City])
				fail[e]++
				failCounter++
				dbg("Drop:[Preload][Postcode][Length]%v %v", row[target.PreLoad.Postcode], row[target.PreLoad.City])
				continue
			}
			t := &tagSet{
				Postcode:  postcode(row[target.PreLoad.Postcode]),
				City:      city(row[target.PreLoad.City]),
				Country:   country(target.Country),
				Preloaded: true,
			}
			if t.uniform() {
				continue
			}
			targets <- t
			dbg("Add:[Preload][%v]%v %v", t.Country, t.Postcode, t.City)
			counter++

		}
		fmt.Printf("\nOSM:PreLoadFile:Fail      # %v", failCounter)
		fmt.Printf("\nOSM:PreLoadFile:Total     # %v", counter)
		writeJsonFile(target.Country, "error.preload.json", fail)
	}
}

// checkPreloadFile ...
func (target *Target) checkPreloadFile() bool {

	// int
	var err error

	// build filname
	target.PreLoad.Filename = filepath.Join(filepath.Dir(target.FileName), "validated-preload", target.Country+".csv")

	// open file handle
	target.PreLoad.File, err = os.Open(target.PreLoad.Filename)
	if err == nil {
		// define preload csv defaults
		target.PreLoad.Fields = 6
		target.PreLoad.City = 2
		target.PreLoad.Postcode = 3
		target.PreLoad.PostcodeLenght = 5
		return true
	}
	if os.IsNotExist(err) {
		target.PreLoad.Filename = "n/a"
		fmt.Printf("\nOSM:PreLoadFile           # %v", target.PreLoad.Filename)
		return false
	}
	target.PreLoad.Filename = err.Error()
	fmt.Printf("\nOSM:PreLoadFile:Error     # %v", target.PreLoad.Filename)
	return false
}
