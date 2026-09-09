package osm2addr

import (
	"os"
	"testing"
)

func TestInitDebug(t *testing.T) {
	cases := []struct {
		env  string
		want bool
	}{
		{"", false},
		{"0", false},
		{"false", false},
		{"no", false},
		{"off", false},
		{"FALSE", false},
		{"  off  ", false},
		{"1", true},
		{"true", true},
		{"yes", true},
		{"on", true},
		{"verbose", true},
		{"TRUE", true},
	}
	for _, c := range cases {
		t.Run(c.env, func(t *testing.T) {
			if c.env == "" {
				if err := os.Unsetenv("DEBUG"); err != nil {
					t.Fatal(err)
				}
			} else {
				t.Setenv("DEBUG", c.env)
			}
			debug = false
			InitDebug()
			if debug != c.want {
				t.Errorf("DEBUG=%q: got %v, want %v", c.env, debug, c.want)
			}
		})
	}
}

func TestInitDebugUnset(t *testing.T) {
	if err := os.Unsetenv("DEBUG"); err != nil {
		t.Fatal(err)
	}
	debug = true
	InitDebug()
	if debug {
		t.Error("debug should be false when DEBUG is unset")
	}
}

func TestDbgNoOpWhenDisabled(t *testing.T) {
	debug = false
	dbg("should not print: %v", 1)
}

func TestDebugEnabled(t *testing.T) {
	if !DebugEnabled() && debug {
		t.Error("DebugEnabled() should track debug var")
	}
}
