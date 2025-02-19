all: 
	make -C cmd/osm2addr all

update-de: 
	mkdir -p data && cd data && curl -O https://download.geofabrik.de/europe/germany-latest.osm.pbf

update-dach:
	mkdir -p data && cd data && curl -O https://download.geofabrik.de/europe/dach-latest.osm.pbf

update-eu: 
	mkdir -p data && cd data && curl -O https://download.geofabrik.de/europe-latest.osm.pbf
