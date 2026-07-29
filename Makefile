PROJECT=$(shell basename $(CURDIR))
PROJECT_PKG=paepcke.de/$(PROJECT)/cmd/$(PROJECT)
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILDTIME=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS=-ldflags "-X $(PROJECT_PKG).version=$(VERSION) -X $(PROJECT_PKG).commit=$(COMMIT) -X $(PROJECT_PKG).buildtime=$(BUILDTIME)"
# CURL=curl --follow 
CURL=aria2c

info:
	echo "osm2addr $(VERSION)"

update: 
	git pull
	git pull --tags

run: build 
	./osm2addr 
	
build:
	go build $(LDFLAGS) -o $(PROJECT) ./cmd/$(PROJECT)

deps: 
	git config core.fileMode false
	rm go.mod go.sum
	go mod init paepcke.de/$(PROJECT)
	go mod tidy -v	

check: 
	gofmt -w -s .
	go vet ./...
	go fix ./...
	CGO_ENABLED=0 staticcheck

##########################
# PROJECT SPECIFIC TASKS #
##########################

update-de: 
	mkdir -p data && $(CURL) -o data/germany-latest.osm.pbf https://download.geofabrik.de/europe/germany-latest.osm.pbf
	# mkdir -p data && $(CURL) -o data/validated && $(CURL) -o data/validated-preload/DE.csv https://downloads.suche-postleitzahl.org/v2/public/zuordnung_plz_ort.csv 

update-dach:
	mkdir -p data && $(CURL) -o data/dach-latest.osm.pbf https://download.geofabrik.de/europe/dach-latest.osm.pbf

update-eu: 
	mkdir -p data && $(CURL) -o data/europe-latest.osm.pbf https://download.geofabrik.de/europe-latest.osm.pbf





