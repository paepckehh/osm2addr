# osm2addr
## Convert OSM OpenStreetMap Open Source Data into Json Address Validation DB Tables (MongoDB,..)
### HOW TO USE 
Generate German (DE) Json mapping tables
```Shell 
curl -O https://download.geofabrik.de/europe/germany-latest.osm.pbf
go run paepcke.de/osm2addr/cmd/osm2addr@latest DE germany-latest.osm.pbf
ls -la json/DE
```
