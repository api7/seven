package DB

import (
	"github.com/hashicorp/go-memdb"
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
)

const (
	Upstream = "Upstream"
)

type UpstreamDB struct {
	Upstreams []*v1.Upstream
}

type UpstreamRequest struct {
	Name string
}

func (ur *UpstreamRequest) FindByName() (*v1.Upstream, error){
	txn := DB.Txn(false)
	defer txn.Abort()
	if raw, err := txn.First(Upstream, "name", ur.Name); err != nil {
		return nil, err
	} else {
		currentUpstream := raw.(*v1.Upstream)
		return currentUpstream, nil
	}
}

// insertUpstream insert upstream to memDB
func (upstreamDB *UpstreamDB) InsertUpstreams() error{
	txn := DB.Txn(true)
	defer txn.Abort()
	for _, u := range upstreamDB.Upstreams{
		if err := txn.Insert(Upstream, u); err != nil {
			return err
		}
	}
	txn.Commit()
	return nil
}

func (upstreamDB *UpstreamDB) UpdateUpstreams() error {
	txn := DB.Txn(true)
	defer txn.Abort()
	for _, u := range upstreamDB.Upstreams {
		// delete
		if _, err := txn.DeleteAll(Upstream, "id", u.ID); err != nil {
			return err
		}
		// insert
		if err := txn.Insert(Upstream, u); err != nil {
			return err
		}
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