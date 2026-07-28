package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

var (
	version   = "dev"
	commit    = "none"
	buildtime = "unknown"
)

func semver() string {
	if version != "dev" && version != "" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" && s.Value != "" {
				return "dev-" + s.Value[:7]
			}
		}
	}
	return "dev"
}

func buildCommit() string {
	if commit != "none" && commit != "" {
		return commit
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" && s.Value != "" {
				return s.Value[:7]
			}
		}
	}
	return "none"
}

func buildTime() string {
	if buildtime != "unknown" && buildtime != "" {
		return buildtime
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			if s.Key == "vcs.time" && s.Value != "" {
				return s.Value
			}
		}
	}
	return "unknown"
}

func versionLine() string {
	return fmt.Sprintf("osm2addr %s (commit %s, built %s, %s/%s)",
		semver(), buildCommit(), buildTime(), runtime.GOOS, runtime.GOARCH)
}
