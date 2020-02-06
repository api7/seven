package DB

import (
	"github.com/hashicorp/go-memdb"
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
)

const (
	Service = "Service"
)

type ServiceDB struct {
	Service *v1.Service
}

func (db *ServiceDB) UpdateService() error{
	txn := DB.Txn(true)
	defer txn.Abort()
	// 1. delete
	if _, err := txn.DeleteAll(Service, "id", db.Service.ID); err != nil {
		return err
	}
	// 2. insert
	if err := txn.Insert(Service, db.Service); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

var serviceSchema = &memdb.TableSchema{
	Name: Service,
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
