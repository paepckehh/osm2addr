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
ls -la DE
```

### EXAMPLE RUN (EXAMPLE DATA OUTPUT see example-json-output folder)
```Shell 
# make
make -C cmd/osm2addr all
make[1]: Entering directory '/home/me/projects/osm2addr/cmd/osm2addr'
GOARCH=amd64 GOAMD64=v3 go run -mod=readonly main.go DE ../../data/germany-latest.osm.pbf

OSM:Startup       # 2025-02-21 09:33:13.140059015 +0000 UTC m=+0.000427720
OSM:TargetCountry # DE
OSM:Worker        # 1
OSM:File          # ../../data/germany-latest.osm.pbf
OSM:File:URL      # https://download.geofabrik.de/europe/germany-updates
OSM:File:Repl:USM # 4330
OSM:File:Repl:TS  # 2025-02-13 21:21:14 +0000 UTC

addr:country:uniq                10  addr:all:country:valid  2.792.983  addr:country:err          0
addr:target:city:uniq         8.045  addr:all:city:valid     3.659.945  addr:city:err            17
addr:target:street:uniq     120.550  addr:all:street:valid   3.905.365  addr:street:err         162
addr:target:postcode:uniq     7.541  addr:all:postcode:valid 3.513.835  addr:postcode:err         0
addr:target:records       2.699.702  addr:all:records        2.711.889  addr:uniform:err         53

Worker#0 Processed => Objects:411.113.874 => Nodes:411.113.874 => AddrTags:71.720.811 => ExitCode:EOF

Total Time Taken: 33.282589674s
```

### FEATURES
- fast: extracts, analyzes and normalizes around 500.000.000 OSM Nodes in under 40 sec (consumer laptop)
