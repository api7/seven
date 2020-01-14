package DB

import "github.com/hashicorp/go-memdb"

const (
	Upstream = "Upstream"
)

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
