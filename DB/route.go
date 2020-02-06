package DB

import (
	"github.com/hashicorp/go-memdb"
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
)

const (
	Route = "Route"
)

type RouteDB struct {
	Route *v1.Route
}

func (db *RouteDB) UpdateRoute() error{
	txn := DB.Txn(true)
	defer txn.Abort()
	// 1. delete
	if _, err := txn.DeleteAll(Route, "id", db.Route.ID); err != nil {
		return err
	}
	// 2. insert
	if err := txn.Insert(Route, db.Route); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

var routeSchema = &memdb.TableSchema{
	Name: Route,
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
