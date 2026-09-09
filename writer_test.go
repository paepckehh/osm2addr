package osm2addr

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteJsonFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeJsonFile(dir, "DE", "test.json", map[string]int{"a": 1})

	b, err := os.ReadFile(filepath.Join(dir, "DE", "test.json"))
	if err != nil {
		t.Fatal(err)
	}
	want := "{\n\t\"a\": 1\n}"
	if string(b) != want {
		t.Errorf("got %q, want %q", b, want)
	}

	info, err := os.Stat(filepath.Join(dir, "DE"))
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Error("country folder must be a directory")
	}
}
