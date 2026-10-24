package decoder

import (
	"bytes"
	"encoding/binary"
	"testing"

	"google.golang.org/protobuf/proto"

	"paepcke.de/osm2addr/internal/protobuf"
)

// headerBlob frames a single blob exactly like an OSM PBF stream: the
// big-endian BlobHeader length, the BlobHeader and the Blob message.
func headerBlob(t *testing.T, blobType string, payload proto.Message) []byte {
	t.Helper()

	data, err := proto.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	blob, err := proto.Marshal(&protobuf.Blob{Raw: data})
	if err != nil {
		t.Fatal(err)
	}
	header, err := proto.Marshal(&protobuf.BlobHeader{
		Type:     proto.String(blobType),
		Datasize: proto.Int32(int32(len(blob))),
	})
	if err != nil {
		t.Fatal(err)
	}

	out := make([]byte, 4, 4+len(header)+len(blob))
	binary.BigEndian.PutUint32(out, uint32(len(header)))
	out = append(out, header...)
	out = append(out, blob...)
	return out
}

func TestLoadHeader(t *testing.T) {
	t.Parallel()

	raw := headerBlob(t, "OSMHeader", &protobuf.HeaderBlock{
		Writingprogram:                   proto.String("osm2addr-test"),
		OsmosisReplicationBaseUrl:        proto.String("https://example.com/"),
		OsmosisReplicationSequenceNumber: proto.Int64(7),
	})

	h, err := LoadHeader(bytes.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	if h.WritingProgram != "osm2addr-test" {
		t.Errorf("writing program: got %q", h.WritingProgram)
	}
	if h.OsmosisReplicationBaseURL != "https://example.com/" || h.OsmosisReplicationSequenceNumber != 7 {
		t.Errorf("unexpected replication metadata: %+v", h)
	}
}

func TestLoadHeaderEmptyDataBlob(t *testing.T) {
	t.Parallel()

	// An OSMData blob without primitive groups yields no elements. LoadHeader
	// must report an error instead of indexing an empty slice.
	raw := headerBlob(t, "OSMData", &protobuf.PrimitiveBlock{
		Stringtable: &protobuf.StringTable{S: []string{""}},
	})

	if _, err := LoadHeader(bytes.NewReader(raw)); err == nil {
		t.Fatal("expected an error for a non-header payload")
	}
}

func TestLoadHeaderTruncated(t *testing.T) {
	t.Parallel()

	if _, err := LoadHeader(bytes.NewReader([]byte{0, 0, 0})); err == nil {
		t.Fatal("expected an error for a truncated header")
	}
}
