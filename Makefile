PROJECT=$(shell basename $(CURDIR))

all:
	make -C cmd/$(PROJECT) all

clean:
	make -C cmd/$(PROJECT) clean

examples:
	make -C cmd/$(PROJECT) examples

deps: 
	rm go.mod go.sum
	go mod init paepcke.de/$(PROJECT)
	go mod tidy -v	

check: 
	gofmt -w -s .
	go vet .
	staticcheck
	golangci-lint run
	make -C cmd/$(PROJECT) check

##########################
# PROJECT SPECIFIC TASKS #
##########################

update-de: 
	mkdir -p data && cd data && curl -O https://download.geofabrik.de/europe/germany-latest.osm.pbf

update-dach:
	mkdir -p data && cd data && curl -O https://download.geofabrik.de/europe/dach-latest.osm.pbf

update-eu: 
	mkdir -p data && cd data && curl -O https://download.geofabrik.de/europe-latest.osm.pbf

