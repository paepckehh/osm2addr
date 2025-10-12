all:
	go run main.go DE ../../data/germany-latest.osm.pbf

clean:
	rm -rf json

check: 
	gofmt -w -s .
	CGO_ENABLED=0 go vet .
	CGO_ENABLED=0 staticcheck
