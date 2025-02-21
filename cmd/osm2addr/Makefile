all:
	GOARCH=amd64 GOAMD64=v3 go run -mod=readonly main.go DE ../../data/germany-latest.osm.pbf

clean:
	rm -rf json

de:
	GOARCH=amd64 GOAMD64=v3 go run -mod=readonly main.go DE ../../data/germany-latest.osm.pbf

examples:
	GOARCH=amd64 GOAMD64=v3 go run -mod=readonly main.go LU ../../data/germany-latest.osm.pbf
	mkdir -p .././example-json-output
	cp -af json/LU ../../example-json-output/

