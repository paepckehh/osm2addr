package osm2addr

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// trustedPreloadFeed ...
func (target *Target) trustedPreloadFeed(targets chan<- *tagSet) {
	if target.checkTrustedPreloadFile() {

		// init
		fail, failCounter, counter := make(map[string]int), 0, 0
		defer func() {
			if err := target.PreLoadTrusted.File.Close(); err != nil {
				fmt.Printf("[OSM2ADDR][ERROR] unable to close file: %s", err)
			}
		}()

		// report
		fmt.Printf("\n----------------------------------------------------------------------------------")
		fmt.Printf("\nOSM:PreLoadTrusted        # %v", target.PreLoadTrusted.Filename)

		// setup
		r := csv.NewReader(target.PreLoadTrusted.File)
		r.Comma = target.PreLoadTrusted.Comma
		r.FieldsPerRecord = -1
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
				e := fmt.Sprintf("[ERROR][PRETRUSTED][CSV][PARSE]#%v#", err)
				fail[e]++
				failCounter++
				dbg("Drop:[PreloadTrusted][CSV][Parse]%v", err)
				continue
			}
			if len(row) < target.PreLoadTrusted.Fields {
				e := fmt.Sprintf("[ERROR][PRETRUSTED][CSV][FIELDS]#%v#", row)
				fail[e]++
				failCounter++
				dbg("Drop:[PreloadTrusted][CSV][Fields]%v", row)
				continue
			}
			t := &tagSet{
				Postcode:  postcode(row[target.PreLoadTrusted.Postcode]),
				City:      city(row[target.PreLoadTrusted.City]),
				Street:    street(row[target.PreLoadTrusted.Street]),
				Country:   country(target.Country),
				Preloaded: true,
			}
			targets <- t
			dbg("Add:[PreloadTrusted][%v]%v %v %v", t.Country, t.Postcode, t.City, t.Street)
			counter++

		}
		fmt.Printf("\nOSM:PreLoadTrusted:Fail   # %v", failCounter)
		fmt.Printf("\nOSM:PreLoadTrusted:Total  # %v", counter)
		writeJsonFile(target.outBase(), target.Country, "error.preload.trusted.json", fail)
	}
}

// checkTrustedPreloadFile ...
func (target *Target) checkTrustedPreloadFile() bool {

	// build filename candidates, pbf dir first, process cwd as fallback
	var candidates []string
	if dir := filepath.Dir(target.FileName); dir != "" && dir != "." {
		candidates = append(candidates, filepath.Join(dir, "preload.csv"))
	}
	candidates = append(candidates, "preload.csv")

	// find first usable preload.csv
	for _, name := range candidates {
		f, err := os.Open(name)
		if err != nil {
			continue
		}
		ok, sep, city, postcode, street := trustedPreloadHeader(f)
		if !ok {
			_ = f.Close()
			target.PreLoadTrusted.Filename = name + " [ERROR][HEADER][SCHEMA][UTF8][SEPARATOR]"
			fmt.Printf("\nOSM:PreLoadTrusted:Error  # %v", target.PreLoadTrusted.Filename)
			return false
		}
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			_ = f.Close()
			continue
		}
		target.PreLoadTrusted.File = f
		target.PreLoadTrusted.Filename = name
		target.PreLoadTrusted.Comma = sep
		target.PreLoadTrusted.City = city
		target.PreLoadTrusted.Postcode = postcode
		target.PreLoadTrusted.Street = street
		target.PreLoadTrusted.Fields = max(city, postcode, street) + 1
		return true
	}

	// report n/a
	target.PreLoadTrusted.Filename = "n/a"
	fmt.Printf("\nOSM:PreLoadTrusted        # %v", target.PreLoadTrusted.Filename)
	return false
}

// trustedPreloadHeader confirms utf-8 encoding, detects the csv separator
// (comma is the assumed default) and locates the exact header field names
// ORT_NAME, POSTLEITZAHL and STRASSE_NAME.
func trustedPreloadHeader(f *os.File) (ok bool, sep rune, city, postcode, street int) {

	// defaults
	sep, city, postcode, street = ',', -1, -1, -1

	// read first line, assume it is a valid csv header
	line, err := bufio.NewReader(f).ReadString('\n')
	if err != nil && line == "" {
		return false, sep, city, postcode, street
	}
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")

	// confirm utf-8 encoding, strip bom
	line = strings.TrimPrefix(line, "\uFEFF")
	if !utf8.ValidString(line) {
		return false, sep, city, postcode, street
	}

	// detect separator, comma is the assumed default
	best := strings.Count(line, string(sep))
	for _, candidate := range []rune{';', '\t', '|'} {
		if c := strings.Count(line, string(candidate)); c > best {
			sep, best = candidate, c
		}
	}

	// locate exact header field names
	for i, name := range strings.Split(line, string(sep)) {
		switch name {
		case "ORT_NAME":
			city = i
		case "POSTLEITZAHL":
			postcode = i
		case "STRASSE_NAME":
			street = i
		}
	}
	if city < 0 || postcode < 0 || street < 0 {
		return false, sep, city, postcode, street
	}
	return true, sep, city, postcode, street
}
