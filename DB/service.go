package DB

import (
	"github.com/hashicorp/go-memdb"
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
	"fmt"
)

const (
	Service = "Service"
)

type ServiceRequest struct {
	Group string
	Name string
}

func (sr *ServiceRequest) FindByName() (*v1.Service, error){
	txn := DB.Txn(false)
	defer txn.Abort()
	if raw, err := txn.First(Service, "name", sr.Group, sr.Name); err != nil {
		return nil, err
	} else {
		if raw != nil {
			currentService := raw.(*v1.Service)
			return currentService, nil
		}
		return nil, fmt.Errorf("NOT FOUND")
	}
}

func (db *ServiceDB) Insert() error {
	txn := DB.Txn(true)
	defer txn.Abort()
	for _, s := range db.Services {
		if err := txn.Insert(Service, s); err != nil {
			return err
		}
	}
	txn.Commit()
	return nil
}

type ServiceDB struct {
	Services []*v1.Service
}

func (db *ServiceDB) UpdateService() error{
	txn := DB.Txn(true)
	defer txn.Abort()
	for _, s := range db.Services {
		// 1. delete
		if _, err := txn.DeleteAll(Service, "name", *(s.Group), *(s.Name)); err != nil {
			return err
		}
		// 2. insert
		if err := txn.Insert(Service, s); err != nil {
			return err
		}
	}

	txn.Commit()
	return nil
}

var serviceSchema = &memdb.TableSchema{
	Name: Service,
	Indexes: map[string]*memdb.IndexSchema{
		"name": {
			Name:         "name",
			Unique:       true,
			Indexer:      indexer(),
			AllowMissing: true,
		},
	},
}
