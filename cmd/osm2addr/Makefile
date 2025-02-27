all:
	go run main.go DE ../../data/germany-latest.osm.pbf

clean:
	rm -rf json

examples:
	go run main.go LU ../../data/germany-latest.osm.pbf
	mkdir -p ../../example-json-output
	cp -af json/LU ../../example-json-output/

check: 
	gofmt -w -s .
	staticcheck
	golangci-lint run
