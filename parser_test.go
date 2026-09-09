package osm2addr

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"google.golang.org/protobuf/proto"
	"paepcke.de/osm2addr/internal/protobuf"
)

// appendPBFBlob writes a single PBF framing unit: the big-endian length of the
// blob header, the BlobHeader and the Blob message wrapping the payload.
func appendPBFBlob(buf *bytes.Buffer, blobType string, payload proto.Message) error {
	data, err := proto.Marshal(payload)
	if err != nil {
		return err
	}
	blob, err := proto.Marshal(&protobuf.Blob{Raw: data})
	if err != nil {
		return err
	}
	header, err := proto.Marshal(&protobuf.BlobHeader{
		Type:     proto.String(blobType),
		Datasize: proto.Int32(int32(len(blob))),
	})
	if err != nil {
		return err
	}
	var size [4]byte
	binary.BigEndian.PutUint32(size[:], uint32(len(header)))
	buf.Write(size[:])
	buf.Write(header)
	buf.Write(blob)
	return nil
}

// writeTestPBF builds a minimal valid OSM PBF file: one OSMHeader blob and one
// OSMData blob with three tagged nodes.
func writeTestPBF(t *testing.T, path string) {
	t.Helper()

	var buf bytes.Buffer
	if err := appendPBFBlob(&buf, "OSMHeader", &protobuf.HeaderBlock{
		Writingprogram:                   proto.String("osm2addr-test"),
		OsmosisReplicationBaseUrl:        proto.String("https://example.com/"),
		OsmosisReplicationSequenceNumber: proto.Int64(1),
	}); err != nil {
		t.Fatal(err)
	}
	if err := appendPBFBlob(&buf, "OSMData", &protobuf.PrimitiveBlock{
		Stringtable: &protobuf.StringTable{S: []string{
			"", "addr:country", "DE", "addr:postcode", "10115",
			"addr:city", "Berlin", "addr:street", "Hauptstrasse",
			"addr:street:name", "AT",
		}},
		Primitivegroup: []*protobuf.PrimitiveGroup{{Nodes: []*protobuf.Node{
			{ // complete address, must be emitted
				Id:   proto.Int64(1),
				Lat:  proto.Int64(1),
				Lon:  proto.Int64(1),
				Keys: []uint32{1, 3, 5, 7},
				Vals: []uint32{2, 4, 6, 8},
			},
			{ // multi-colon addr tag must not be misread as addr:street
				Id:   proto.Int64(2),
				Lat:  proto.Int64(1),
				Lon:  proto.Int64(1),
				Keys: []uint32{1, 9},
				Vals: []uint32{2, 8},
			},
			{ // country mismatch, must be dropped
				Id:   proto.Int64(3),
				Lat:  proto.Int64(1),
				Lon:  proto.Int64(1),
				Keys: []uint32{1, 3, 5, 7},
				Vals: []uint32{10, 4, 6, 8},
			},
		}}},
	}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestPBFParser(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	pbfPath := filepath.Join(dir, "test.osm.pbf")
	writeTestPBF(t, pbfPath)

	f, err := os.Open(pbfPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = f.Close()
	}()

	target := &Target{Country: "DE", FileName: pbfPath, File: f, outDir: dir}
	targets := make(chan *tagSet, 8)
	done := make(chan struct{})
	go func() {
		defer close(done)
		pbfparser(target, targets)
	}()
	<-done
	close(targets)

	var got []*tagSet
	for ts := range targets {
		got = append(got, ts)
	}
	if len(got) != 1 {
		t.Fatalf("expected exactly one complete address, got %v", got)
	}
	want := &tagSet{Country: "DE", Postcode: "10115", City: "Berlin", Street: "Hauptstraße"}
	if *got[0] != *want {
		t.Errorf("got %+v, want %+v", *got[0], *want)
	}
}
