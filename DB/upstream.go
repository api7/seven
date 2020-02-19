package DB

import (
	"fmt"
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
	"github.com/hashicorp/go-memdb"
)

const (
	Upstream = "Upstream"
)

type UpstreamDB struct {
	Upstreams []*v1.Upstream
}

type UpstreamRequest struct {
	Group string
	Name  string
}

func (ur *UpstreamRequest) FindByName() (*v1.Upstream, error) {
	txn := DB.Txn(false)
	defer txn.Abort()
	if raw, err := txn.First(Upstream, "name", ur.Group, ur.Name); err != nil {
		return nil, err
	} else {
		if raw != nil {
			currentUpstream := raw.(*v1.Upstream)
			return currentUpstream, nil
		}
		return nil, fmt.Errorf("NOT FOUND")
	}
}

// insertUpstream insert upstream to memDB
func (upstreamDB *UpstreamDB) InsertUpstreams() error {
	txn := DB.Txn(true)
	defer txn.Abort()
	for _, u := range upstreamDB.Upstreams {
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
		if _, err := txn.DeleteAll(Upstream, "name", *(u.Group), *(u.Name)); err != nil {
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
		"name": {
			Name:         "name",
			Unique:       true,
			Indexer:      indexer(),
			AllowMissing: true,
		},
	},
}

func indexer() *memdb.CompoundIndex{
	var idx = make([]memdb.Indexer, 0)
	idx = append(idx, &memdb.StringFieldIndex{Field: "Group"})
	idx = append(idx, &memdb.StringFieldIndex{Field: "Name"})
	return &memdb.CompoundIndex{Indexes: idx, AllowMissing: false}
}

