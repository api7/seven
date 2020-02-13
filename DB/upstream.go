package DB

import (
	"github.com/hashicorp/go-memdb"
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
)

const (
	Upstream = "Upstream"
)

type UpstreamDB struct {
	Upstream *v1.Upstream
}

func (upstreamDB UpstreamDB) UpdateUpstream() error {
	txn := DB.Txn(true)
	defer txn.Abort()
	// delete
	if _, err := txn.DeleteAll(Upstream, "id", upstreamDB.Upstream.ID); err != nil {
		return err
	}
	// insert
	if err := txn.Insert(Upstream, upstreamDB.Upstream); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

var upstreamSchema = &memdb.TableSchema{
	Name: Upstream,
	Indexes: map[string]*memdb.IndexSchema{
		"id": {
			Name:    "id",
			Unique:  true,
			Indexer: &memdb.StringFieldIndex{Field: "ID"},
		},
		"name": {
			Name:         "name",
			Unique:       true,
			Indexer:      &memdb.StringFieldIndex{Field: "Name"},
			AllowMissing: true,
		},
	},
}