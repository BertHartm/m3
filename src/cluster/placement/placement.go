// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package placement

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errInvalidHostShards   = errors.New("invalid shards assigned to a host")
	errDuplicatedShards    = errors.New("invalid placement, there are duplicated shards in one replica")
	errUnexpectedShards    = errors.New("invalid placement, there are unexpected shard ids on host")
	errTotalShardsMismatch = errors.New("invalid placement, the total shards in the placement does not match expected number")
	errInvalidShardsCount  = errors.New("invalid placement, the count for a shard does not match replica factor")
)

// snapshot implements Snapshot
type snapshot struct {
	hostShards []HostShards
	rf         int
	shards     []uint32
}

// NewPlacementSnapshot returns a placement snapshot
func NewPlacementSnapshot(hss []HostShards, ids []uint32, rf int) Snapshot {
	return snapshot{hostShards: hss, rf: rf, shards: ids}
}

// NewEmptyPlacementSnapshot returns an empty placement
func NewEmptyPlacementSnapshot(hosts []Host, ids []uint32) Snapshot {
	hostShards := make([]HostShards, len(hosts), len(hosts))
	for i, ph := range hosts {
		hostShards[i] = NewEmptyHostShards(ph)
	}

	return snapshot{hostShards: hostShards, shards: ids, rf: 0}
}

func (ps snapshot) HostShards() []HostShards {
	result := make([]HostShards, ps.HostsLen())
	for i, hs := range ps.hostShards {
		result[i] = hs
	}
	return result
}

func (ps snapshot) HostsLen() int {
	return len(ps.hostShards)
}

func (ps snapshot) Replicas() int {
	return ps.rf
}

func (ps snapshot) ShardsLen() int {
	return len(ps.shards)
}

func (ps snapshot) Shards() []uint32 {
	return ps.shards
}

func (ps snapshot) HostShard(id string) HostShards {
	for _, phs := range ps.HostShards() {
		if phs.Host().ID() == id {
			return phs
		}
	}
	return nil
}

func (ps snapshot) Validate() error {
	shardCountMap := ConvertShardSliceToMap(ps.shards)
	if len(shardCountMap) != len(ps.shards) {
		return errDuplicatedShards
	}

	expectedTotal := len(ps.shards) * ps.rf
	actualTotal := 0
	for _, hs := range ps.hostShards {
		for _, id := range hs.Shards() {
			if count, exist := shardCountMap[id]; exist {
				shardCountMap[id] = count + 1
				continue
			}

			return errUnexpectedShards
		}
		actualTotal += hs.ShardsLen()
	}

	if expectedTotal != actualTotal {
		return errTotalShardsMismatch
	}

	for shard, c := range shardCountMap {
		if ps.rf != c {
			return fmt.Errorf("invalid shard count for shard %d: expected %d, actual %d", shard, ps.rf, c)
		}
	}
	return nil
}

// Copy copies a snapshot
func (ps snapshot) Copy() Snapshot {
	return snapshot{hostShards: copyHostShards(ps.HostShards()), rf: ps.Replicas(), shards: ps.Shards()}
}

func copyHostShards(hss []HostShards) []HostShards {
	copied := make([]HostShards, len(hss))
	for i, hs := range hss {
		copied[i] = NewHostShards(hs.Host(), hs.Shards())
	}
	return copied
}

// ConvertShardSliceToMap is an util function that converts a slice of shards to a map
func ConvertShardSliceToMap(ids []uint32) map[uint32]int {
	shardCounts := make(map[uint32]int)
	for _, id := range ids {
		shardCounts[id] = 0
	}
	return shardCounts
}

// hostShards implements HostShards
type hostShards struct {
	host      Host
	shardsSet map[uint32]struct{}
}

// NewHostShards returns a HostShards with shards
func NewHostShards(host Host, shards []uint32) HostShards {
	m := make(map[uint32]struct{}, len(shards))
	for _, s := range shards {
		m[s] = struct{}{}
	}
	return &hostShards{host: host, shardsSet: m}
}

// NewEmptyHostShards returns a HostShards with no shards assigned
func NewEmptyHostShards(host Host) HostShards {
	return NewHostShards(host, nil)
}

func (h hostShards) Host() Host {
	return h.host
}

func (h hostShards) Shards() []uint32 {
	s := make([]uint32, 0, len(h.shardsSet))
	for shard := range h.shardsSet {
		s = append(s, shard)
	}
	return s
}

func (h hostShards) AddShard(s uint32) {
	h.shardsSet[s] = struct{}{}
}

func (h hostShards) RemoveShard(shard uint32) {
	delete(h.shardsSet, shard)
}

func (h hostShards) ContainsShard(shard uint32) bool {
	if _, exist := h.shardsSet[shard]; exist {
		return true
	}
	return false
}

func (h hostShards) ShardsLen() int {
	return len(h.shardsSet)
}

// NewHost returns a Host
func NewHost(id, rack, zone string, weight uint32) Host {
	return host{id: id, rack: rack, zone: zone, weight: weight}
}

type host struct {
	id     string
	rack   string
	zone   string
	weight uint32
}

func (h host) ID() string {
	return h.id
}

func (h host) Rack() string {
	return h.rack
}

func (h host) Zone() string {
	return h.zone
}

func (h host) Weight() uint32 {
	return h.weight
}

func (h host) String() string {
	return fmt.Sprintf("[id:%s, rack:%s, zone:%s, weight:%v]", h.id, h.rack, h.zone, h.weight)
}

// ByIDAscending sorts Hosts by ID
type ByIDAscending []Host

func (s ByIDAscending) Len() int {
	return len(s)
}

func (s ByIDAscending) Less(i, j int) bool {
	return strings.Compare(s[i].ID(), s[j].ID()) < 0
}

func (s ByIDAscending) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
