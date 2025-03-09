package osm2addr

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

// preloadFeed ...
func (target *Target) preloadFeed() {
	if target.checkPreloadFile() {
		fmt.Printf("\nOSM:PreLoadFile       #  %v", target.PreLoad.Filename)
		defer target.PreLoad.File.Close()
		counter := 0
		r := csv.NewReader(target.PreLoad.File)
		r.Comma = ','
		r.FieldsPerRecord = target.PreLoad.Fields
		r.LazyQuotes = false
		r.ReuseRecord = true
		for {
			row, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			if len(row[target.PreLoad.City]) < 3 {
				// log.Printf("PreLoad:Error:CityLenght %v \t => %v ", row[target.PreLoad.Postcode], row[target.PreLoad.City])
				continue
			}
			if len(row[target.PreLoad.Postcode]) != target.PreLoad.PostcodeLenght {
				// log.Printf("PreLoad:Error:PostcodeLenght %v \t => %v ", row[target.PreLoad.Postcode], row[target.PreLoad.City])
				continue
			}
			_, err = strconv.Atoi(row[target.PreLoad.Postcode])
			if err != nil {
				// log.Printf("PreLoad:Error:IsInteger %v \t => %v ", row[target.PreLoad.Postcode], row[target.PreLoad.City])
				continue
			}
			targets <- &TagSET{
				Postcode: postcode(row[target.PreLoad.Postcode]),
				City:     city(row[target.PreLoad.City]),
			}
			counter++

		}
		fmt.Printf("\nOSM:PreLoadFile:Done  #  %v", counter)
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
		fmt.Printf("\nOSM:PreLoadFile       #  %v", target.PreLoad.Filename)
		return false
	}
	target.PreLoad.Filename = err.Error()
	fmt.Printf("\nOSM:PreLoadFile:Error #  %v", target.PreLoad.Filename)
	return false
}
