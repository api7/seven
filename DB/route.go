package DB

import (
	"github.com/hashicorp/go-memdb"
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
	"fmt"
)

const (
	Route = "Route"
)

type RouteRequest struct {
	Name string
}

func (rr *RouteRequest) FindByName() (*v1.Route, error){
	txn := DB.Txn(false)
	defer txn.Abort()
	if raw, err := txn.First(Route, "name", rr.Name); err != nil {
		return nil, err
	} else {
		if raw != nil {
			currentRoute := raw.(*v1.Route)
			return currentRoute, nil
		}
		return nil, fmt.Errorf("NOT FOUND")
	}
}

type RouteDB struct {
	Routes []*v1.Route
}

// InsertRoute insert route to memDB
func (db *RouteDB) Insert() error{
	txn := DB.Txn(true)
	defer txn.Abort()
	for _, r := range db.Routes {
		if err := txn.Insert(Route, r); err != nil {
			return err
		}
	}
	txn.Commit()
	return nil
}

func (db *RouteDB) UpdateRoute() error{
	txn := DB.Txn(true)
	defer txn.Abort()
	for _, r := range db.Routes {
		// 1. delete
		if _, err := txn.DeleteAll(Route, "id", r.ID); err != nil {
			return err
		}
		// 2. insert
		if err := txn.Insert(Route, r); err != nil {
			return err
		}
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
