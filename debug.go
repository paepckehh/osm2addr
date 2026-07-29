package osm2addr

import (
	"fmt"
	"os"
	"strings"
)

var debug bool

// InitDebug reads the DEBUG environment variable at program start. If DEBUG
// is set to any non-empty value (except "0", "false", "no", "off"), verbose
// event logging is enabled.
func InitDebug() {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("DEBUG")))
	switch v {
	case "", "0", "false", "no", "off":
		debug = false
	default:
		debug = true
	}
}

// DebugEnabled reports whether verbose event logging is active.
func DebugEnabled() bool { return debug }

// dbg emits a verbose event line when DEBUG is enabled. The format follows
// the existing OSM:<Section>:<Key> # <value> convention so debug lines mix
// naturally with the standard stat output.
func dbg(format string, args ...any) {
	if !debug {
		return
	}
	fmt.Printf("\nOSM:DEBUG:"+format, args...)
}
