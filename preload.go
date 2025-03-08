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
		defer target.PreLoadFile.Close()
		r := csv.NewReader(target.PreLoadFile)
		for {
			line, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			_ = line
		}
	}
	fmt.Printf("\nOSM:PreLoadFile       #  %v", target.PreLoadFilename)
	preload.Done()
}

// checkPreloadFile ...
func (target *Target) checkPreloadFile() bool {

	// int
	var err error

	// build filname
	target.PreLoadFilename = filepath.Join(filepath.Dir(target.FileName), "validated-preload", target.Country+".csv")

	// open file handle
	target.PreLoadFile, err = os.Open(target.PreLoadFilename)
	if err == nil {
		return false
	}
	if os.IsNotExist(err) {
		target.PreLoadFilename = "n/a"
		return false
	}
	target.PreLoadFilename = err.Error()
	return true
}
