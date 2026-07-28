PROJECT=$(shell basename $(CURDIR))
PROJECT_PKG=paepcke.de/$(PROJECT)/cmd/$(PROJECT)
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILDTIME=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS=-ldflags "-X $(PROJECT_PKG).version=$(VERSION) -X $(PROJECT_PKG).commit=$(COMMIT) -X $(PROJECT_PKG).buildtime=$(BUILDTIME)"

all:
	go run $(LDFLAGS) main.go DE ../../data/germany-latest.osm.pbf

build:
	go build $(LDFLAGS) -o $(PROJECT) .

clean:
	rm -rf json

check: 
	gofmt -w -s .
	go vet ./...
	go fix ./...
	CGO_ENABLED=0 staticcheck
