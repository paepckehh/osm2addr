package osm2addr

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	pbfPath := filepath.Join(dir, "test.osm.pbf")
	writeTestPBF(t, pbfPath)

	trusted := "POSTLEITZAHL;ORT_NAME;STRASSE_NAME\n" +
		"10115;Berlin;Torstraße\n"
	if err := os.WriteFile(filepath.Join(dir, "preload.csv"), []byte(trusted), 0644); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(pbfPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = f.Close()
	}()

	target := &Target{Country: "DE", FileName: pbfPath, File: f, outDir: dir}
	if err := Parse(target); err != nil {
		t.Fatal(err)
	}

	var addr map[string]map[string]map[string]bool
	readJson(t, filepath.Join(dir, "DE", "addr.json"), &addr)
	streets := addr["10115"]["Berlin"]
	if len(streets) != 2 {
		t.Fatalf("expected two Berlin streets (trusted preload + pbf), got %v", streets)
	}
	if !streets["Torstraße"] {
		t.Error("trusted preload street missing")
	}
	if !streets["Hauptstraße"] {
		t.Error("pbf street missing")
	}

	var places map[string]map[string]map[string]string
	readJson(t, filepath.Join(dir, "DE", "addr2placeID.json"), &places)
	if len(places["10115"]["Berlin"]) != 2 {
		t.Errorf("expected two place IDs for Berlin, got %v", places["10115"]["Berlin"])
	}

	if _, err := os.Stat(filepath.Join(dir, "DE", "placeID2addr.json")); err != nil {
		t.Error("placeID2addr.json must be written")
	}
	if _, err := os.Stat(filepath.Join(dir, "DE", "warning.json")); err != nil {
		t.Error("warning.json must be written")
	}
	if _, err := os.Stat(filepath.Join(dir, "DE", "corrected.json")); err != nil {
		t.Error("corrected.json must be written")
	}
	if _, err := os.Stat(filepath.Join(dir, "DE", "error.preload.trusted.json")); err != nil {
		t.Error("error.preload.trusted.json must be written")
	}
}
