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
	"time"

	"google.golang.org/protobuf/proto"

	"paepcke.de/osm2addr/internal/model"
	"paepcke.de/osm2addr/internal/protobuf"
)

func parsePrimitiveBlock(buffer []byte) ([]model.Object, error) {
	pb := &protobuf.PrimitiveBlock{}
	if err := proto.Unmarshal(buffer, pb); err != nil {
		return nil, err
	}
	c := newBlockContext(pb)
	elements := make([]model.Object, 0)
	for _, pg := range pb.GetPrimitivegroup() {
		elements = append(elements, c.decodeNodes(pg.GetNodes())...)
		elements = append(elements, c.decodeDenseNodes(pg.GetDense())...)
	}
	return elements, nil
}

func (c *blockContext) decodeNodes(nodes []*protobuf.Node) (elements []model.Object) {
	elements = make([]model.Object, len(nodes))
	for i, node := range nodes {
		elements[i] = &model.Node{
			ID:   model.ID(node.GetId()),
			Tags: c.decodeTags(node.GetKeys(), node.GetVals()),
			Info: c.decodeInfo(node.GetInfo()),
		}
	}
	return elements
}

func (c *blockContext) decodeDenseNodes(nodes *protobuf.DenseNodes) []model.Object {
	ids := nodes.GetId()
	elements := make([]model.Object, len(ids))
	tic := c.newTagsContext(nodes.GetKeysVals())
	dic := c.newDenseInfoContext(nodes.GetDenseinfo())
	var id int64
	for i := range ids {
		id = ids[i] + id
		elements[i] = &model.Node{
			ID:   model.ID(id),
			Tags: tic.decodeTags(),
			Info: dic.decodeInfo(i),
		}
	}
	return elements
}

func (c *blockContext) decodeTags(keyIDs, valIDs []uint32) map[string]string {
	tags := make(map[string]string, len(keyIDs))
	if len(keyIDs) > len(valIDs) {
		return tags
	}
	for i, keyID := range keyIDs {
		tags[c.strings[keyID]] = c.strings[valIDs[i]]
	}
	return tags
}

func (c *blockContext) decodeInfo(info *protobuf.Info) *model.Info {
	i := &model.Info{Visible: true}
	if info != nil {
		i.Version = info.GetVersion()
		i.Timestamp = toTimestamp(c.dateGranularity, int64(info.GetTimestamp()))
		i.Changeset = info.GetChangeset()
		i.UID = model.UID(info.GetUid())

		i.User = c.strings[info.GetUserSid()]

		if info.Visible != nil {
			i.Visible = info.GetVisible()
		}
	}
	return i
}

type blockContext struct {
	strings         []string
	granularity     int32
	dateGranularity int32
}

func newBlockContext(pb *protobuf.PrimitiveBlock) *blockContext {
	return &blockContext{
		strings:         pb.GetStringtable().GetS(),
		granularity:     pb.GetGranularity(),
		dateGranularity: pb.GetDateGranularity(),
	}
}

type denseInfoContext struct {
	version         int32
	timestamp       int64
	changeset       int64
	uid             model.UID
	userSid         int32
	dateGranularity int32
	strings         []string
	versions        []int32
	uids            []model.UID
	timestamps      []int64
	changesets      []int64
	userSids        []int32
	visibilities    []bool
}

func (c *blockContext) newDenseInfoContext(di *protobuf.DenseInfo) *denseInfoContext {
	uids := make([]model.UID, len(di.GetUid()))
	for i, uid := range di.GetUid() {
		uids[i] = model.UID(uid)
	}
	return &denseInfoContext{
		dateGranularity: c.dateGranularity,
		strings:         c.strings,
		versions:        di.GetVersion(),
		uids:            uids,
		timestamps:      di.GetTimestamp(),
		changesets:      di.GetChangeset(),
		userSids:        di.GetUserSid(),
		visibilities:    di.GetVisible(),
	}
}

func (dic *denseInfoContext) decodeInfo(i int) *model.Info {
	info := &model.Info{Visible: true}
	if i < len(dic.versions) {
		dic.version += dic.versions[i]
		info.Version = dic.version
	}
	if i < len(dic.uids) {
		dic.uid += dic.uids[i]
		info.UID = dic.uid
	}
	if i < len(dic.timestamps) {
		dic.timestamp += dic.timestamps[i]
		info.Timestamp = toTimestamp(dic.dateGranularity, dic.timestamp)
	}
	if i < len(dic.changesets) {
		dic.changeset += dic.changesets[i]
		info.Changeset = dic.changeset
	}
	if i < len(dic.userSids) {
		dic.userSid += dic.userSids[i]
		info.User = dic.strings[dic.userSid]
	}
	if i < len(dic.visibilities) {
		info.Visible = dic.visibilities[i]
	}
	return info
}

type tagsContext struct {
	strings []string
	i       int
	keyVals []int32
}

func (c *blockContext) newTagsContext(keyVals []int32) *tagsContext {
	return &tagsContext{strings: c.strings, keyVals: keyVals}
}

func (tic *tagsContext) decodeTags() map[string]string {
	tags := map[string]string{}
	i := tic.i
	for i+1 < len(tic.keyVals) && tic.keyVals[i] > 0 {
		tags[tic.strings[tic.keyVals[i]]] = tic.strings[tic.keyVals[i+1]]
		i += 2
	}
	tic.i = i + 1
	return tags
}

// toTimestamp converts a timestamp with a specific granularity, in units of
// milliseconds, to a UTC timestamp of type Time.
func toTimestamp(granularity int32, timestamp int64) time.Time {
	ms := time.Duration(timestamp) * time.Duration(granularity) * time.Millisecond
	return time.Unix(0, int64(ms)).UTC()
}
