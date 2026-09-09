package osm2addr

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPreloadFeedMalformedRow(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "validated-preload"), 0755); err != nil {
		t.Fatal(err)
	}
	csvData := "f0,f1,city,postcode,f4,f5\n" +
		"10115,Berlin\n" +
		"a,b,Berlin,10115,x,y\n"
	if err := os.WriteFile(filepath.Join(dir, "validated-preload", "DE.csv"), []byte(csvData), 0644); err != nil {
		t.Fatal(err)
	}

	target := &Target{Country: "DE", FileName: filepath.Join(dir, "germany.pbf"), outDir: dir}
	if !target.checkPreloadFile() {
		t.Fatal("preload file must be detected")
	}

	targets := make(chan *tagSet)
	done := make(chan struct{})
	go func() {
		defer close(done)
		target.preloadFeed(targets)
	}()

	select {
	case ts := <-targets:
		if !ts.Preloaded {
			t.Error("entry must be marked preloaded")
		}
		if ts.City != "Berlin" || string(ts.Postcode) != "10115" {
			t.Errorf("unexpected entry: city %q postcode %q", ts.City, ts.Postcode)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("expected one preloaded entry")
	}
	<-done

	select {
	case <-targets:
		t.Error("unexpected extra entry")
	default:
	}

	b, err := os.ReadFile(filepath.Join(dir, "DE", "error.preload.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fail map[string]int
	if err := json.Unmarshal(b, &fail); err != nil {
		t.Fatal(err)
	}
	if len(fail) != 1 {
		t.Errorf("expected exactly one failure case, got %v", fail)
	}
}

func TestPreloadFeedValidation(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "validated-preload"), 0755); err != nil {
		t.Fatal(err)
	}
	csvData := "f0,f1,city,postcode,f4,f5\n" +
		"a,b,Berlin,10115,x,y\n" + // ok
		"a,b,B,10115,x,y\n" + // city too short
		"a,b,Berlin,1011,x,y\n" + // postcode too short
		"a,b,Berlin,101155,x,y\n" // postcode too long
	if err := os.WriteFile(filepath.Join(dir, "validated-preload", "DE.csv"), []byte(csvData), 0644); err != nil {
		t.Fatal(err)
	}

	target := &Target{Country: "DE", FileName: filepath.Join(dir, "germany.pbf"), outDir: dir}
	if !target.checkPreloadFile() {
		t.Fatal("preload file must be detected")
	}

	targets := make(chan *tagSet, 4)
	target.preloadFeed(targets)
	close(targets)

	var got []*tagSet
	for ts := range targets {
		got = append(got, ts)
	}
	if len(got) != 1 {
		t.Fatalf("expected exactly one valid entry, got %v", got)
	}
	if got[0].City != "Berlin" || string(got[0].Postcode) != "10115" {
		t.Errorf("unexpected entry: %+v", got[0])
	}

	b, err := os.ReadFile(filepath.Join(dir, "DE", "error.preload.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fail map[string]int
	if err := json.Unmarshal(b, &fail); err != nil {
		t.Fatal(err)
	}
	if len(fail) != 3 {
		t.Errorf("expected three failure cases, got %v", fail)
	}
}

func TestCheckPreloadFileMissing(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	target := &Target{Country: "DE", FileName: filepath.Join(dir, "germany.pbf")}
	if target.checkPreloadFile() {
		t.Error("missing preload file must be skipped")
	}
	if target.PreLoad.Filename != "n/a" {
		t.Errorf("unexpected filename: %v", target.PreLoad.Filename)
	}
}
