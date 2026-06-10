all:
	go run main.go DE ../../data/germany-latest.osm.pbf

clean:
	rm -rf json

check: 
	gofmt -w -s .
	go vet ./...
	go fix ./...
	CGO_ENABLED=0 staticcheck
