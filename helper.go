package osm2addr

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
)

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

func writeJsonFile(filename string, inMap map[string]map[string]bool) {
	folder := "json"
	_ = os.Mkdir(folder, 0755)
	j, err := json.Marshal(inMap)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(folder, filename), j, 0666); err != nil {
		panic(err)
	}
}

func tryNormaliseGermanCity(content string) string {
	content = strings.ReplaceAll(content, " a.d. ", " an der ")
	content = strings.ReplaceAll(content, " a.d.", " an der ")
	content = strings.ReplaceAll(content, " i. ", " im ")
	content = strings.ReplaceAll(content, " i.", " im ")
	content = strings.ReplaceAll(content, " a. ", " am ")
	content = strings.ReplaceAll(content, " a.", " am ")
	content = strings.ReplaceAll(content, " b. ", " bei ")
	content = strings.ReplaceAll(content, " b.", " bei ")
	content = strings.ReplaceAll(content, " v.d. ", " von der ")
	content = strings.ReplaceAll(content, " v. d. ", " von der ")
	content = strings.ReplaceAll(content, " v.d.", " von der ")
	content = strings.ReplaceAll(content, " v. d.", " von der ")
	return content
}
