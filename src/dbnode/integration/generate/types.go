	"github.com/m3db/m3db/src/dbnode/clock"
	"github.com/m3db/m3db/src/dbnode/encoding"
	"github.com/m3db/m3db/src/dbnode/sharding"
	"github.com/m3db/m3db/src/dbnode/ts"
	// Write writes the data
	Write(ns ident.ID, shards sharding.ShardSet, data SeriesBlocksByStart) error