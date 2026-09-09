package osm2addr

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// writeJsonFile ...
func writeJsonFile(base, countrycode, filename string, in any) {
	folder := filepath.Join(base, countrycode)
	if err := os.MkdirAll(folder, 0755); err != nil {
		panic(err)
	}
	j, err := json.MarshalIndent(&in, "", "\t")
	if err != nil {
		panic(err)
	}
	file := filepath.Join(folder, filename)
	if err := os.WriteFile(file, j, 0644); err != nil {
		panic(err)
	}
	fmt.Printf("\nOSM:Writer:JSON           # %v", file)
}
