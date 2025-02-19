# osm2addr
## Convert OSM OpenStreetMap Open Source Data into Json Address Validation DB Tables (MongoDB,..)
## https://download.geofabrik.de
### HOW TO USE (Example OSM DATA GERMANY) 
```Shell 
curl -O https://download.geofabrik.de/europe/germany-latest.osm.pbf
go run paepcke.de/osm2addr/cmd/osm2addr@latest DE germany-latest.osm.pbf
ls -la json/DE
```
