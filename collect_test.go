package osm2addr

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCollect(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	target := &Target{Country: "DE", outDir: dir}
	entries := []*tagSet{
		{Country: "DE", Postcode: "10115", City: "Berlin", Street: "Hauptstraße", Preloaded: true},
		{Country: "DE", Postcode: "10115", City: "Berlin", Street: "Hauptstraße", Preloaded: true}, // dedup
		{Country: "DE", Postcode: "10115", City: "Berlin", Street: "Torstraße"},                    // added street
		{Country: "DE", Postcode: "10115", City: "Berlino", Street: "Sonnenallee"},                 // levenshtein warning
		{Country: "DE", Postcode: "65123", City: "Bad-Homburg", Street: "Gartenweg", Preloaded: true},
		{Country: "DE", Postcode: "65123", City: "Bad Homburg", Street: "Lindenweg"}, // separator correction
		{Country: "DE", Postcode: "20095", City: "Hamburg", Street: "", Preloaded: true},
		{Country: "DE", Postcode: "20095", City: "Hamburg", Street: "Mönckebergstraße"},
		{Country: "DE", Postcode: "99999", City: "Zwickau", Street: "Hauptplatz"}, // non-preloaded postcode warning
	}

	targets := make(chan *tagSet, len(entries))
	for _, e := range entries {
		targets <- e
	}
	close(targets)
	collect(target, targets)

	var addr map[string]map[string]map[string]bool
	readJson(t, filepath.Join(dir, "DE", "addr.json"), &addr)
	if len(addr["10115"]["Berlin"]) != 2 {
		t.Errorf("expected two Berlin streets, got %v", addr["10115"]["Berlin"])
	}
	if len(addr["10115"]["Berlino"]) != 1 {
		t.Errorf("expected one Berlino street, got %v", addr["10115"]["Berlino"])
	}
	if len(addr["65123"]["Bad-Homburg"]) != 2 {
		t.Errorf("expected two corrected Bad-Homburg streets, got %v", addr["65123"]["Bad-Homburg"])
	}
	if _, ok := addr["65123"]["Bad Homburg"]; ok {
		t.Error("separator variant city must be merged, not exported")
	}
	if len(addr["20095"]["Hamburg"]) != 1 {
		t.Errorf("expected one Hamburg street, got %v", addr["20095"]["Hamburg"])
	}
	if len(addr["99999"]["Zwickau"]) != 1 {
		t.Errorf("expected one Zwickau street, got %v", addr["99999"]["Zwickau"])
	}

	var places map[string]map[string]map[string]string
	readJson(t, filepath.Join(dir, "DE", "addr2placeID.json"), &places)
	for _, city := range places["65123"] {
		for _, pid := range city {
			if len(pid) != 24 {
				t.Errorf("place id %q must be 24 hex chars", pid)
			}
		}
	}

	var p map[string]tagSet
	readJson(t, filepath.Join(dir, "DE", "placeID2addr.json"), &p)
	if len(p) != 7 {
		t.Errorf("expected seven unique places, got %d", len(p))
	}

	var corrected map[string]int
	readJson(t, filepath.Join(dir, "DE", "corrected.json"), &corrected)
	if len(corrected) != 1 {
		t.Errorf("expected one corrected case, got %v", corrected)
	}

	var warning map[string]int
	readJson(t, filepath.Join(dir, "DE", "warning.json"), &warning)
	if len(warning) != 2 {
		t.Errorf("expected two warning cases, got %v", warning)
	}
}

func TestCollectPreloadedCitySeededWithoutStreet(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	target := &Target{Country: "DE", outDir: dir}
	entries := []*tagSet{
		{Country: "DE", Postcode: "20095", City: "Hamburg", Street: "", Preloaded: true},
		{Country: "DE", Postcode: "20095", City: "Hamburg", Street: "", Preloaded: true}, // dedup, no placeID
	}

	targets := make(chan *tagSet, len(entries))
	for _, e := range entries {
		targets <- e
	}
	close(targets)
	collect(target, targets)

	var addr map[string]map[string]map[string]bool
	readJson(t, filepath.Join(dir, "DE", "addr.json"), &addr)
	if _, ok := addr["20095"]["Hamburg"]; !ok {
		t.Error("city-only entry must seed the city mapping")
	}
	if len(addr["20095"]["Hamburg"]) != 0 {
		t.Errorf("city-only entry must not seed streets, got %v", addr["20095"]["Hamburg"])
	}

	var p map[string]tagSet
	readJson(t, filepath.Join(dir, "DE", "placeID2addr.json"), &p)
	if len(p) != 0 {
		t.Errorf("street-less entries must not produce place IDs, got %v", p)
	}
}

func readJson(t *testing.T, path string, v any) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, v); err != nil {
		t.Fatal(err)
	}
}
