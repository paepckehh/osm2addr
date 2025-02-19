package osm2addr

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
)

// hu print large number readable for humans and fixed lenght
func hu(in int) string {
	h := humanize.Comma(int64(in))
	for {
		if len(h) < 10 {
			h = " " + h
			continue
		}
		return h
	}

}

// writeJsonFile ...
func writeJsonFile(countrycode, filename string, inMap map[string]map[string]bool) {
	folder := filepath.Join("json", countrycode)
	_ = os.MkdirAll(folder, 0755)
	j, err := json.Marshal(inMap)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(folder, filename), j, 0644); err != nil {
		panic(err)
	}
}
