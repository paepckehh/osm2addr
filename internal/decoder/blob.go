// Copyright 2017-25 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package decoder

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/klauspost/compress/zlib"

	"github.com/destel/rill"
	"google.golang.org/protobuf/proto"
	"paepcke.de/osm2addr/internal/core"
	"paepcke.de/osm2addr/internal/model"
	"paepcke.de/osm2addr/internal/protobuf"
)

type blob struct {
	header *protobuf.BlobHeader
	blob   *protobuf.Blob
}

func Generate(ctx context.Context, reader io.Reader) func(yield func(enc blob, err error) bool) {
	return func(yield func(enc blob, err error) bool) {
		buffer := core.NewPooledBuffer()
		defer func() {
			if err := buffer.Close(); err != nil {
				fmt.Printf("[OSM2ADDR][ERROR] close: %v", err)
			}
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			h, err := readBlobHeader(buffer, reader)
			if err != nil {
				if err != io.EOF {
					slog.Error(err.Error())
					yield(blob{}, err)
				}
				return
			}

			b, err := readBlob(buffer, reader, h)
			if err != nil {
				slog.Error(err.Error())
				yield(blob{}, err)

				return
			}

			if !yield(blob{header: h, blob: b}, nil) {
				return
			}

			buffer.Reset()
		}
	}
}

func Decode(array []blob) (out <-chan rill.Try[[]model.Object]) {
	ch := make(chan rill.Try[[]model.Object])
	out = ch

	go func() {
		defer close(ch)

		for _, enc := range array {
			elements, err := extract(enc.header, enc.blob)
			if err != nil {
				slog.Error(err.Error())
				ch <- rill.Try[[]model.Object]{Error: err}

				return
			}
			ch <- rill.Try[[]model.Object]{Value: elements}
		}
	}()
	return
}

// readBlobHeader unmarshals a header from an array of protobuf encoded bytes.
// The header is used when decoding blobs into OSM elements.
func readBlobHeader(buffer *core.PooledBuffer, rdr io.Reader) (header *protobuf.BlobHeader, err error) {
	var size uint32
	err = binary.Read(rdr, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}
	buffer.Reset()
	if _, err := io.CopyN(buffer, rdr, int64(size)); err != nil {
		return nil, err
	}
	header = &protobuf.BlobHeader{}
	if err := proto.Unmarshal(buffer.Bytes(), header); err != nil {
		return nil, err
	}
	return header, nil
}

// readBlob unmarshals a blob from an array of protobuf encoded bytes.  The
// blob still needs to be decoded into OSM elements using decode().
func readBlob(buffer *core.PooledBuffer, rdr io.Reader, header *protobuf.BlobHeader) (*protobuf.Blob, error) {
	size := header.GetDatasize()
	buffer.Reset()
	if _, err := io.CopyN(buffer, rdr, int64(size)); err != nil {
		return nil, err
	}
	blob := &protobuf.Blob{}
	if err := proto.Unmarshal(buffer.Bytes(), blob); err != nil {
		return nil, err
	}
	return blob, nil
}

// elements unmarshals an array of OSM elements from an array of protobuf encoded
// bytes.  The bytes could possibly be compressed; zlibBuf is used to facilitate
// decompression.
func extract(header *protobuf.BlobHeader, blob *protobuf.Blob) ([]model.Object, error) {
	var buf []byte

	switch {
	case blob.Raw != nil:
		buf = blob.GetRaw()

	case blob.ZlibData != nil:
		zlibBuf := core.NewPooledBuffer()
		defer func() {
			if err := zlibBuf.Close(); err != nil {
				fmt.Printf("[OSM2ADDR][ERROR] close :%v", err)
			}
		}()
		r, err := zlib.NewReader(bytes.NewReader(blob.GetZlibData()))
		if err != nil {
			return nil, err
		}
		zlibBuf.Reset()
		rawBufferSize := int(blob.GetRawSize() + bytes.MinRead)
		if rawBufferSize > zlibBuf.Cap() {
			zlibBuf.Grow(rawBufferSize)
		}
		_, err = zlibBuf.ReadFrom(r)
		if err != nil {
			return nil, err
		}
		if zlibBuf.Len() != int(blob.GetRawSize()) {
			err = fmt.Errorf("raw blob data size %d but expected %d", zlibBuf.Len(), blob.GetRawSize())

			return nil, err
		}
		buf = zlibBuf.Bytes()
	default:
		return nil, errors.New("unknown blob data type")
	}
	ht := *header.Type
	switch ht {
	case "OSMHeader":
		{
			h, err := parseOSMHeader(buf)
			if err != nil {
				return nil, err
			}

			return []model.Object{h}, nil
		}
	case "OSMData":
		return parsePrimitiveBlock(buf)
	default:
		return nil, fmt.Errorf("unknown header type %s", ht)
	}
}

// parseOSMHeader unmarshals the OSM header from an array of protobuf encoded bytes.
func parseOSMHeader(buffer []byte) (*model.Header, error) {
	hb := &protobuf.HeaderBlock{}
	if err := proto.Unmarshal(buffer, hb); err != nil {
		return nil, err
	}
	header := &model.Header{
		RequiredFeatures:                 hb.GetRequiredFeatures(),
		OptionalFeatures:                 hb.GetOptionalFeatures(),
		WritingProgram:                   hb.GetWritingprogram(),
		Source:                           hb.GetSource(),
		OsmosisReplicationBaseURL:        hb.GetOsmosisReplicationBaseUrl(),
		OsmosisReplicationSequenceNumber: hb.GetOsmosisReplicationSequenceNumber(),
	}

	if hb.OsmosisReplicationTimestamp != nil {
		header.OsmosisReplicationTimestamp = time.Unix(*hb.OsmosisReplicationTimestamp, 0)
	}

	return header, nil
}
