package osm2addr

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTrustedPreloadFeed(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	csvData := "POSTLEITZAHL;ORT_NAME;STRASSE_NAME;LAND\n" +
		"1234;Berlin;Hauptstrasse;DE\n" +
		"10115;Berlin;Chausseestrasse;DE\n"
	if err := os.WriteFile(filepath.Join(dir, "preload.csv"), []byte(csvData), 0644); err != nil {
		t.Fatal(err)
	}

	target := &Target{Country: "DE", FileName: filepath.Join(dir, "germany.pbf"), outDir: dir}
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

	targets := make(chan *tagSet)
	done := make(chan struct{})
	go func() {
		defer close(done)
		target.trustedPreloadFeed(targets)
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
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "preload.csv"), []byte("a,b,c\n1,2,3\n"), 0644); err != nil {
		t.Fatal(err)
	}
	target := &Target{Country: "DE", FileName: filepath.Join(dir, "germany.pbf")}
	if target.checkTrustedPreloadFile() {
		t.Error("file without required header fields must be rejected")
	}
}

func TestTrustedPreloadFeedMissing(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	target := &Target{Country: "DE", FileName: filepath.Join(dir, "germany.pbf")}
	if target.checkTrustedPreloadFile() {
		t.Error("missing preload.csv must be skipped")
	}
	if target.PreLoadTrusted.Filename != "n/a" {
		t.Errorf("unexpected filename: %v", target.PreLoadTrusted.Filename)
	}
}

func TestTrustedPreloadHeader(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		csv      string
		ok       bool
		sep      rune
		city     int
		postcode int
		street   int
	}{
		{
			name:     "semicolon separated",
			csv:      "POSTLEITZAHL;ORT_NAME;STRASSE_NAME\n1;2;3\n",
			ok:       true,
			sep:      ';',
			postcode: 0,
			city:     1,
			street:   2,
		},
		{
			name:     "tab separated with bom",
			csv:      "\uFEFFORT_NAME\tPOSTLEITZAHL\tSTRASSE_NAME\n1\t2\t3\n",
			ok:       true,
			sep:      '\t',
			city:     0,
			postcode: 1,
			street:   2,
		},
		{
			name:     "comma separated default",
			csv:      "ORT_NAME,POSTLEITZAHL,STRASSE_NAME\n1,2,3\n",
			ok:       true,
			sep:      ',',
			city:     0,
			postcode: 1,
			street:   2,
		},
		{
			name: "missing street column",
			csv:  "POSTLEITZAHL;ORT_NAME\n1;2\n",
			ok:   false,
		},
		{
			name: "empty file",
			csv:  "",
			ok:   false,
		},
		{
			name: "non utf8 header",
			csv:  "POSTLEITZAHL;ORT_NAME;STRASSE_NAME;\xff\xfe\n",
			ok:   false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			f, err := os.CreateTemp(t.TempDir(), "preload*.csv")
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = f.Close()
			}()
			if _, err := f.WriteString(c.csv); err != nil {
				t.Fatal(err)
			}
			if _, err := f.Seek(0, 0); err != nil {
				t.Fatal(err)
			}
			ok, sep, city, postcode, street := trustedPreloadHeader(f)
			if ok != c.ok {
				t.Fatalf("ok: got %v, want %v", ok, c.ok)
			}
			if !c.ok {
				return
			}
			if sep != c.sep || city != c.city || postcode != c.postcode || street != c.street {
				t.Errorf("got sep %q city %v postcode %v street %v, want %q %v %v %v",
					sep, city, postcode, street, c.sep, c.city, c.postcode, c.street)
			}
		})
	}
}
