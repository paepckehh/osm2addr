all: 
	make -C cmd/osm2addr all

clean:
	make -C cmd/osm2addr clean

examples:
	make -C cmd/osm2addr examples

update-de: 
	mkdir -p data && cd data && curl -O https://download.geofabrik.de/europe/germany-latest.osm.pbf

update-dach:
	mkdir -p data && cd data && curl -O https://download.geofabrik.de/europe/dach-latest.osm.pbf

update-eu: 
	mkdir -p data && cd data && curl -O https://download.geofabrik.de/europe-latest.osm.pbf

deps: 
	rm go.mod go.sum
	go mod init paepcke.de/osm2addr
	go mod tidy -v 
