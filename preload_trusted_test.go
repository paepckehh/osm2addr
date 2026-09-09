package osm2addr

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTrustedPreloadFeed(t *testing.T) {
	t.Chdir(t.TempDir())

	dir := t.TempDir()
	csvData := "POSTLEITZAHL;ORT_NAME;STRASSE_NAME;LAND\n" +
		"1234;Berlin;Hauptstrasse;DE\n" +
		"10115;Berlin;Chausseestrasse;DE\n"
	if err := os.WriteFile(filepath.Join(dir, "preload.csv"), []byte(csvData), 0644); err != nil {
		t.Fatal(err)
	}

	target := &Target{Country: "DE", FileName: filepath.Join(dir, "germany.pbf")}
	if !target.checkTrustedPreloadFile() {
		t.Fatal("trusted preload file must be detected")
	}
	if target.PreLoadTrusted.Comma != ';' {
		t.Errorf("expected ';' separator, got %q", target.PreLoadTrusted.Comma)
	}
	if target.PreLoadTrusted.Postcode != 0 || target.PreLoadTrusted.City != 1 || target.PreLoadTrusted.Street != 2 {
		t.Errorf("unexpected column indices: postcode %v city %v street %v",
			target.PreLoadTrusted.Postcode, target.PreLoadTrusted.City, target.PreLoadTrusted.Street)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		target.trustedPreloadFeed()
	}()

	var got []*tagSet
	for range 2 {
		select {
		case ts := <-targets:
			if !ts.Preloaded {
				t.Error("entry must be marked preloaded")
			}
			got = append(got, ts)
		case <-time.After(2 * time.Second):
			t.Fatal("expected two trusted entries")
		}
	}
	<-done

	if got[0].City != "Berlin" || string(got[0].Postcode) != "1234" || got[0].Street != "Hauptstrasse" {
		t.Errorf("unexpected first entry: %+v", got[0])
	}
	if got[1].City != "Berlin" || string(got[1].Postcode) != "10115" || got[1].Street != "Chausseestrasse" {
		t.Errorf("unexpected second entry: %+v", got[1])
	}

	select {
	case <-targets:
		t.Error("unexpected extra entry")
	default:
	}
}

func TestTrustedPreloadFeedInvalidHeader(t *testing.T) {
	t.Chdir(t.TempDir())

	if err := os.WriteFile("preload.csv", []byte("a,b,c\n1,2,3\n"), 0644); err != nil {
		t.Fatal(err)
	}
	target := &Target{Country: "DE", FileName: "germany.pbf"}
	if target.checkTrustedPreloadFile() {
		t.Error("file without required header fields must be rejected")
	}
}

func TestTrustedPreloadFeedMissing(t *testing.T) {
	t.Chdir(t.TempDir())

	target := &Target{Country: "DE", FileName: "germany.pbf"}
	if target.checkTrustedPreloadFile() {
		t.Error("missing preload.csv must be skipped")
	}
	if target.PreLoadTrusted.Filename != "n/a" {
		t.Errorf("unexpected filename: %v", target.PreLoadTrusted.Filename)
	}
}
