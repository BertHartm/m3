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
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnapshot(t *testing.T) {
	h1 := NewEmptyHostShards(NewHost("r1h1", "r1", "z1", 1))
	h1.AddShard(1)
	h1.AddShard(2)
	h1.AddShard(3)

	h2 := NewEmptyHostShards(NewHost("r2h2", "r2", "z1", 1))
	h2.AddShard(4)
	h2.AddShard(5)
	h2.AddShard(6)

	h3 := NewEmptyHostShards(NewHost("r3h3", "r3", "z1", 1))
	h3.AddShard(1)
	h3.AddShard(3)
	h3.AddShard(5)

	h4 := NewEmptyHostShards(NewHost("r4h4", "r4", "z1", 1))
	h4.AddShard(2)
	h4.AddShard(4)
	h4.AddShard(6)

	h5 := NewEmptyHostShards(NewHost("r5h5", "r5", "z1", 1))
	h5.AddShard(5)
	h5.AddShard(6)
	h5.AddShard(1)

	h6 := NewEmptyHostShards(NewHost("r6h6", "r6", "z1", 1))
	h6.AddShard(2)
	h6.AddShard(3)
	h6.AddShard(4)

	hss := []HostShards{h1, h2, h3, h4, h5, h6}

	ids := []uint32{1, 2, 3, 4, 5, 6}
	s := NewPlacementSnapshot(hss, ids, 3)
	assert.NoError(t, s.Validate())

	hs := s.HostShard("r6h6")
	assert.Equal(t, h6, hs)
	hs = s.HostShard("h100")
	assert.Nil(t, hs)

	assert.Equal(t, 6, s.HostsLen())
	assert.Equal(t, 3, s.Replicas())
	assert.Equal(t, ids, s.Shards())
	assert.Equal(t, hss, s.HostShards())

	s = NewEmptyPlacementSnapshot([]Host{NewHost("h1", "r1", "z1", 1), NewHost("h2", "r2", "z1", 1)}, ids)
	assert.Equal(t, 0, s.Replicas())
	assert.Equal(t, ids, s.Shards())
	assert.NoError(t, s.Validate())
}

func TestValidate(t *testing.T) {
	ids := []uint32{1, 2, 3, 4, 5, 6}

	h1 := NewEmptyHostShards(NewHost("r1h1", "r1", "z1", 1))
	h1.AddShard(1)
	h1.AddShard(2)
	h1.AddShard(3)

	h2 := NewEmptyHostShards(NewHost("r2h2", "r2", "z1", 1))
	h2.AddShard(4)
	h2.AddShard(5)
	h2.AddShard(6)

	hss := []HostShards{h1, h2}
	s := NewPlacementSnapshot(hss, ids, 1)
	assert.NoError(t, s.Validate())

	// mismatch shards
	s = NewPlacementSnapshot(hss, append(ids, 7), 1)
	assert.Error(t, s.Validate())
	assert.Error(t, s.Validate())

	// host missing a shard
	h1 = NewEmptyHostShards(NewHost("r1h1", "r1", "z1", 1))
	h1.AddShard(1)
	h1.AddShard(2)
	h1.AddShard(3)
	h1.AddShard(4)
	h1.AddShard(5)
	h1.AddShard(6)

	h2 = NewEmptyHostShards(NewHost("r2h2", "r2", "z1", 1))
	h2.AddShard(2)
	h2.AddShard(3)
	h2.AddShard(4)
	h2.AddShard(5)
	h2.AddShard(6)

	hss = []HostShards{h1, h2}
	s = NewPlacementSnapshot(hss, ids, 2)
	assert.Error(t, s.Validate())
	assert.Equal(t, errTotalShardsMismatch, s.Validate())

	// host contains shard that's unexpected to be in snapshot
	h1 = NewEmptyHostShards(NewHost("r1h1", "r1", "z1", 1))
	h1.AddShard(1)
	h1.AddShard(2)
	h1.AddShard(3)
	h1.AddShard(4)
	h1.AddShard(5)
	h1.AddShard(6)
	h1.AddShard(7)

	h2 = NewEmptyHostShards(NewHost("r2h2", "r2", "z1", 1))
	h2.AddShard(2)
	h2.AddShard(3)
	h2.AddShard(4)
	h2.AddShard(5)
	h2.AddShard(6)

	hss = []HostShards{h1, h2}
	s = NewPlacementSnapshot(hss, ids, 2)
	assert.Error(t, s.Validate())
	assert.Equal(t, errUnexpectedShards, s.Validate())

	// duplicated shards
	h1 = NewEmptyHostShards(NewHost("r1h1", "r1", "z1", 1))
	h1.AddShard(2)
	h1.AddShard(3)
	h1.AddShard(4)

	h2 = NewEmptyHostShards(NewHost("r2h2", "r2", "z1", 1))
	h2.AddShard(4)
	h2.AddShard(5)
	h2.AddShard(6)

	hss = []HostShards{h1, h2}
	s = NewPlacementSnapshot(hss, []uint32{2, 3, 4, 4, 5, 6}, 1)
	assert.Error(t, s.Validate())
	assert.Equal(t, errDuplicatedShards, s.Validate())

	// three shard 2 and only one shard 4
	h1 = NewEmptyHostShards(NewHost("r1h1", "r1", "z1", 1))
	h1.AddShard(1)
	h1.AddShard(2)
	h1.AddShard(3)

	h2 = NewEmptyHostShards(NewHost("r2h2", "r2", "z1", 1))
	h2.AddShard(2)
	h2.AddShard(3)
	h2.AddShard(4)

	h3 := NewEmptyHostShards(NewHost("r3h3", "r3", "z1", 1))
	h3.AddShard(1)
	h3.AddShard(2)

	hss = []HostShards{h1, h2, h3}
	s = NewPlacementSnapshot(hss, []uint32{1, 2, 3, 4}, 2)
	assert.Error(t, s.Validate())
}

func TestHostShards(t *testing.T) {
	h1 := NewEmptyHostShards(NewHost("r1h1", "r1", "z1", 1))
	h1.AddShard(1)
	h1.AddShard(2)
	h1.AddShard(3)

	assert.Equal(t, "[id:r1h1, rack:r1, zone:z1, weight:1]", h1.Host().String())

	assert.True(t, h1.ContainsShard(1))
	assert.False(t, h1.ContainsShard(100))
	assert.Equal(t, 3, h1.ShardsLen())
	assert.Equal(t, "r1h1", h1.Host().ID())
	assert.Equal(t, "r1", h1.Host().Rack())

	h1.RemoveShard(1)
	assert.False(t, h1.ContainsShard(1))
	assert.False(t, h1.ContainsShard(100))
	assert.Equal(t, 2, h1.ShardsLen())
	assert.Equal(t, "r1h1", h1.Host().ID())
	assert.Equal(t, "r1", h1.Host().Rack())
}

func TestCopy(t *testing.T) {
	h1 := NewEmptyHostShards(NewHost("r1h1", "r1", "z1", 1))
	h1.AddShard(1)
	h1.AddShard(2)
	h1.AddShard(3)

	h2 := NewEmptyHostShards(NewHost("r2h2", "r2", "z1", 1))
	h2.AddShard(4)
	h2.AddShard(5)
	h2.AddShard(6)

	hss := []HostShards{h1, h2}

	ids := []uint32{1, 2, 3, 4, 5, 6}
	s := NewPlacementSnapshot(hss, ids, 1)
	copy := s.Copy()
	assert.Equal(t, s.HostsLen(), copy.HostsLen())
	assert.Equal(t, s.Shards(), copy.Shards())
	assert.Equal(t, s.Replicas(), copy.Replicas())
	for _, hs := range s.HostShards() {
		assert.Equal(t, copy.HostShard(hs.Host().ID()), hs)
		// make sure they are different objects, updating one won't update the other
		hs.AddShard(100)
		assert.NotEqual(t, copy.HostShard(hs.Host().ID()), hs)
	}
}

func TestSortHostByID(t *testing.T) {
	h1 := NewHost("h1", "", "", 1)
	h2 := NewHost("h2", "", "", 1)
	h3 := NewHost("h3", "", "", 1)
	h4 := NewHost("h4", "", "", 1)
	h5 := NewHost("h5", "", "", 1)
	h6 := NewHost("h6", "", "", 1)

	hs := []Host{h1, h6, h4, h2, h3, h5}
	sort.Sort(ByIDAscending(hs))

	assert.Equal(t, []Host{h1, h2, h3, h4, h5, h6}, hs)
}

func TestOptions(t *testing.T) {
	o := NewOptions()
	assert.False(t, o.LooseRackCheck())
	assert.False(t, o.AllowPartialReplace())
	o = o.SetLooseRackCheck(true)
	assert.True(t, o.LooseRackCheck())
	o = o.SetAllowPartialReplace(true)
	assert.True(t, o.AllowPartialReplace())

	dopts := NewDeploymentOptions()
	assert.Equal(t, defaultMaxStepSize, dopts.MaxStepSize())
	dopts = dopts.SetMaxStepSize(5)
	assert.Equal(t, 5, dopts.MaxStepSize())
}
