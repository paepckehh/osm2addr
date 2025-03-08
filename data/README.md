# OVERVIEW 
[![Go Reference](https://pkg.go.dev/badge/paepcke.de/opnborg.svg)](https://pkg.go.dev/paepcke.de/opnborg) 
[![Go Report Card](https://goreportcard.com/badge/paepcke.de/opnborg)](https://goreportcard.com/report/paepcke.de/opnborg) 
[![Go Build](https://github.com/paepckehh/opnborg/actions/workflows/golang.yml/badge.svg)](https://github.com/paepckehh/opnborg/actions/workflows/golang.yml)
[![License](https://img.shields.io/github/license/paepckehh/opnborg)](https://github.com/paepckehh/opnborg/blob/master/LICENSE)
[![SemVer](https://img.shields.io/github/v/release/paepckehh/opnborg)](https://github.com/paepckehh/opnborg/releases/latest)
<br>[![built with nix](https://builtwithnix.org/badge.svg)](https://search.nixos.org/packages?channel=unstable&from=0&size=50&sort=relevance&type=packages&query=opnborg)

[paepcke.de/opnborg](https://paepcke.de/opnborg/)


# osm2addr
## Convert OSM OpenStreetMap Open Source Data into Json Address Validation DB Tables (MongoDB,..)
### HOW TO USE 
Generate German (DE) Json mapping tables
```Shell 
curl --output germany-latest.osm.pbf https://download.geofabrik.de/europe/germany-latest.osm.pbf
go run paepcke.de/osm2addr/cmd/osm2addr@latest DE germany-latest.osm.pbf
ls -la json/DE
```

Download optional German (DE) validated-preload dataset for postcode:location mapping
```Shell 
curl --output validated-preload/DE.csv https://downloads.suche-postleitzahl.org/v2/public/zuordnung_plz_ort.csv 
```
