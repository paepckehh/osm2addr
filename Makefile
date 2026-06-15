PROJECT=$(shell basename $(CURDIR))
# CURL=curl --follow 
CURL=aria2c

all:
	make -C cmd/$(PROJECT) all

clean:
	make -C cmd/$(PROJECT) clean

deps: 
	rm go.mod go.sum
	go mod init paepcke.de/$(PROJECT)
	go mod tidy -v	

check: 
	gofmt -w -s .
	go vet ./...
	go fix ./...
	CGO_ENABLED=0 staticcheck
	make -C cmd/$(PROJECT) check

##########################
# PROJECT SPECIFIC TASKS #
##########################

update-de: 
	mkdir -p data && $(CURL) -o data/germany-latest.osm.pbf https://download.geofabrik.de/europe/germany-latest.osm.pbf
	# mkdir -p data/validated && $(CURL) -o data/validated-preload/DE.csv https://downloads.suche-postleitzahl.org/v2/public/zuordnung_plz_ort.csv 

update-dach:
	mkdir -p data && $(CURL) -o data/dach-latest.osm.pbf https://download.geofabrik.de/europe/dach-latest.osm.pbf

update-eu: 
	mkdir -p data && $(CURL) -o data/europe-latest.osm.pbf https://download.geofabrik.de/europe-latest.osm.pbf

