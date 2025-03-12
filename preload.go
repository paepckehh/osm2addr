package osm2addr

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// preloadFeed ...
func (target *Target) preloadFeed() {
	if target.checkPreloadFile() {

		// init
		fail, failCounter, counter, uniformCounter, ok := make(map[string]int), 0, 0, 0, false
		defer target.PreLoad.File.Close()

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
				if err == io.EOF {
					break
				}
				panic(err)
			}
			if len(row[target.PreLoad.City]) < 2 {
				e := fmt.Sprintf("[ERROR][PRE][CITY][LENGHT]#%v#%v#", row[target.PreLoad.Postcode], row[target.PreLoad.City])
				if _, ok = fail[e]; !ok {
					fail[e] = 0
				}
				fail[e]++
				failCounter++
				continue
			}
			if len(row[target.PreLoad.Postcode]) != target.PreLoad.PostcodeLenght {
				e := fmt.Sprintf("[ERROR][PRE][POSTCODE][LENGHT]#%v#%v#", row[target.PreLoad.Postcode], row[target.PreLoad.City])
				if _, ok = fail[e]; !ok {
					fail[e] = 0
				}
				fail[e]++
				failCounter++
				continue
			}
			t := &tagSet{
				Postcode: postcode(row[target.PreLoad.Postcode]),
				City:     city(row[target.PreLoad.City]),
				Country:  country(target.Country),
			}
			uniformCounter += t.uniform()
			targets <- t
			counter++

		}
		fmt.Printf("\nOSM:PreLoadFile:Fail      # %v", failCounter)
		fmt.Printf("\nOSM:PreLoadFile:Total     # %v", counter)
		writeJsonFile(target.Country, "error.preload.json", fail)
		preload.Done()
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
