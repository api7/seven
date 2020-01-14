package DB

import "github.com/hashicorp/go-memdb"

var DB *memdb.MemDB

func init(){
	if db, err := NewDB(); err != nil {
		panic(err)
	}else {
		DB = db
	}
}

func NewDB() (*memdb.MemDB, error){
	var schema = &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			Service:  serviceSchema,
			Route:    routeSchema,
			Upstream: upstreamSchema,
		},
	}

	if memDB, err := memdb.NewMemDB(schema); err != nil {
		return nil , err
	} else {
		return memDB, nil
	}
}
