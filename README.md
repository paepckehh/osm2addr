# OVERVIEW 
[![Go Reference](https://pkg.go.dev/badge/paepcke.de/osm2addr.svg)](https://pkg.go.dev/paepcke.de/osm2addr) 
[![Go Report Card](https://goreportcard.com/badge/paepcke.de/osm2addr)](https://goreportcard.com/report/paepcke.de/osm2addr) 
[![Go Build](https://github.com/paepckehh/osm2addr/actions/workflows/golang.yml/badge.svg)](https://github.com/paepckehh/osm2addr/actions/workflows/golang.yml)
[![License](https://img.shields.io/github/license/paepckehh/osm2addr)](https://github.com/paepckehh/osm2addr/blob/master/LICENSE)
[![SemVer](https://img.shields.io/github/v/release/paepckehh/osm2addr)](https://github.com/paepckehh/osm2addr/releases/latest)
<br>[![built with nix](https://builtwithnix.org/badge.svg)](https://search.nixos.org/packages?channel=unstable&from=0&size=50&sort=relevance&type=packages&query=osm2addr)

[paepcke.de/osm2addr](https://paepcke.de/osm2addr/)


# osm2addr
## Parse, extract and uniform OSM OpenStreetMap Open Source Data into json (postal) address validation DB mapping tables (MongoDB,..)
### HOW TO USE
Generate German (DE) Json mapping tables
requirements: curl golang 1.23+ 
```Shell 
curl -O https://download.geofabrik.de/europe/germany-latest.osm.pbf
go run paepcke.de/osm2addr/cmd/osm2addr@latest DE germany-latest.osm.pbf
ls -la json/DE
```
### FEATURES
- fast: extracts, analyzes and normalizes around 500.000.000 OSM Nodes / minute on mid-range consumer hardware
