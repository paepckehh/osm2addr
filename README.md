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

```Shell
mkdir -p data/validated-preload
curl --output data/germany-latest.osm.pbf https://download.geofabrik.de/europe/germany-latest.osm.pbf 
curl --output data/validated-preload/DE.csv https://downloads.suche-postleitzahl.org/v2/public/zuordnung_plz_ort.csv 
go run paepcke.de/osm2addr/cmd/osm2addr@latest DE data/germany-latest.osm.pbf

OSM:Startup               # 2025-03-10 09:23:45.921159027 +0000 UTC m=+0.000386754
OSM:TargetCountry         # DE
OSM:WorkerScale           # 1
OSM:File                  # germany-latest.osm.pbf
----------------------------------------------------------------------------------
OSM:PreLoadFile           # validated-preload/DE.csv
OSM:PreLoadFile:Fail      # 1
OSM:PreLoadFile:Total     # 12853
OSM:Writer:JSON           # json/DE/error.preload.json
----------------------------------------------------------------------------------
OSM:PBF:File:URL          # https://download.geofabrik.de/europe/germany-updates
OSM:PBF:File:Repl:USM     # 4330
OSM:PBF:File:Repl:TS      # 2025-02-13 21:21:14 +0000 UTC
----------------------------------------------------------------------------------
OSM:PBF:Parsed:Objects    # 411.113.874
OSM:PBF:Parsed:Tags       #  71.720.811
OSM:PBF:Parsed:Country    #   2.792.983
OSM:PBF:Parsed:Street     #   3.905.365
OSM:PBF:Parsed:City       #   3.659.945
OSM:PBF:Parsed:Postcode   #   3.513.835
----------------------------------------------------------------------------------
OSM:PBF:Complete:AddrTags #   2.711.889
OSM:PBF:Uniq:Country      #          10
OSM:PBF:Err:Uniform       #      21.854
OSM:PBF:Err:Country       #           0
OSM:PBF:Err:Postcode      #           0
OSM:PBF:Err:City          #          17
OSM:PBF:Err:Street        #         162
----------------------------------------------------------------------------------
OSM:Corrected:Auto:Cases  #          78
OSM:Corrected:Auto:Total  #      21.082
OSM:Corrected:Warn:Cases  #           8
OSM:Corrected:Warn:Total  #           8
OSM:Collect:Places:Total  #     278.381
----------------------------------------------------------------------------------
OSM:Writer:JSON           # json/DE/id.json
OSM:Writer:JSON           # json/DE/warning.json
OSM:Writer:JSON           # json/DE/corrected.json
OSM:Time:Total            # 45.329429572s
```

### FEATURES
- fast: extracts, analyzes and normalizes around 400.000.000 OSM Nodes in under 50 sec (consumer laptop)
