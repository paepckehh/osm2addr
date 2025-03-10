package osm2addr

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// writeJsonFile ...
func writeJsonFile(countrycode, filename string, in interface{}) {
	folder := filepath.Join("json", countrycode)
	_ = os.MkdirAll(folder, 0755)
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
