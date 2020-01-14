package apisix

import (
	"github.com/gxthrj/seven/DB"
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
)

// insertUpstream insert upstream to memDB
func InsertUpstreams(upstreams []*v1.Upstream) error{
	txn := DB.DB.Txn(true)
	for _, u := range upstreams{
		if err := txn.Insert(DB.Upstream, u); err != nil {
			return err
		}
	}
	txn.Commit()
	return nil
}

// InsertRoute insert route to memDB
func InsertRoute(routes []*v1.Route) error{
	txn := DB.DB.Txn(true)
	for _, r := range routes {
		if err := txn.Insert(DB.Route, r); err != nil {
			return err
		}
	}
	txn.Commit()
	return nil
}