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
requirements: curl/wget and golang 1.23+ 

```Shell 
curl -O https://download.geofabrik.de/europe/germany-latest.osm.pbf
go run paepcke.de/osm2addr/cmd/osm2addr@latest DE germany-latest.osm.pbf
ls -la DE
```

### EXAMPLE RUN (EXAMPLE DATA OUTPUT see example-json-output folder)
```Shell
go run main.go DE ../../data/germany-latest.osm.pbf

OSM:Startup           #  2025-03-03 09:03:09.632247303 +0000 UTC m=+0.000512489
OSM:TargetCountry     #  DE
OSM:WorkerScale       #  1
OSM:File              #  ../../data/germany-latest.osm.pbf
OSM:PBF:File:URL      #  https://download.geofabrik.de/europe/germany-updates
OSM:PBF:File:Repl:USM #  4330
OSM:PBF:File:Repl:TS  #  2025-02-13 21:21:14 +0000 UTC
[INFO] Levenshtein:1 # Postcode:17398 # City:Bugewitz <===> City:Bugewtz
[INFO] Levenshtein:1 # Postcode:04779 # City:Wermsdorf <===> City:Wermsorf
[INFO] Levenshtein:1 # Postcode:94051 # City:Hauzenberg <===> City:Haunzenberg
[INFO] Levenshtein:1 # Postcode:22395 # City:Hamburg <===> City:Hanburg
OSM:PBF:ObjectsParsed # 411.113.874
OSM:PBF:AddrTags      #   2.711.889
OSM:PBF:Uniq:Country  #          10
OSM:PBF:Err:Uniform   #      21.858
OSM:PBF:Err:Country   #           0
OSM:PBF:Err:Postcode  #           0
OSM:PBF:Err:City      #          17
OSM:PBF:Err:Street    #         162
OSM:Uniq:City         #       8.036
OSM:Uniq:Street       #     120.548
OSM:Uniq:Postcode     #       7.541
OSM:Collect:Sets      #     265.636
OSM:TotalTime:        #  50.25799055s
```

### FEATURES
- fast: extracts, analyzes and normalizes around 500.000.000 OSM Nodes in under 40 sec (consumer laptop)
