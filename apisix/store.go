package apisix

import (
	"github.com/gxthrj/seven/DB"
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
	"fmt"
	"github.com/gxthrj/seven/conf"
)

// insertUpstream insert upstream to memDB
func InsertUpstreams(upstreams []*v1.Upstream) error{
	txn := DB.DB.Txn(true)
	defer txn.Abort()
	for _, u := range upstreams{
		if err := txn.Insert(DB.Upstream, u); err != nil {
			return err
		}
	}
	txn.Commit()
	return nil
}

func InsertServices(services []*v1.Service) error {
	txn := DB.DB.Txn(true)
	defer txn.Abort()
	for _, s := range services {
		if err := txn.Insert(DB.Service, s); err != nil {
			return err
		}
	}
	txn.Commit()
	return nil
}


// InsertRoute insert route to memDB
func InsertRoute(routes []*v1.Route) error{
	txn := DB.DB.Txn(true)
	defer txn.Abort()
	for _, r := range routes {
		if err := txn.Insert(DB.Route, r); err != nil {
			return err
		}
	}
	txn.Commit()
	return nil
}
// FindRoute find current route in memDB
func FindRoute(route *v1.Route) (*v1.Route,error){
	txn := DB.DB.Txn(false)
	defer txn.Abort()
	raw, _ := txn.First(DB.Route, "name", route.Name)
	if raw != nil { // update
		currentRoute := raw.(*v1.Route)
		return currentRoute, nil
	} else {
		return nil, fmt.Errorf("NOT FOUND")
	}
}
// FindUpstreamByName find upstream from memDB,
// if Not Found, find upstream from apisix
func FindUpstreamByName(name string) (*v1.Upstream, error){
	txn := DB.DB.Txn(false)
	defer txn.Abort()
	raw, _ := txn.First(DB.Upstream, "name", name)
	if raw != nil {
		currentUpstream := raw.(*v1.Upstream)
		return currentUpstream, nil
	} else {
		// find upstream from apisix
		if upstreams, err := ListUpstream(conf.BaseUrl); err != nil {
			// todo log error
		}else {
			for _, upstream := range upstreams {
				if *(upstream.Name) == name {
					// and save to memDB
					InsertUpstreams([]*v1.Upstream{upstream})
					// return
					return upstream, nil
				}
			}
		}

	}
	return nil, nil
}

func FindServiceByName(name string) (*v1.Service, error){
	txn := DB.DB.Txn(false)
	defer txn.Abort()
	raw, _ := txn.First(DB.Service, "name", name)
	if raw != nil {
		currentService := raw.(*v1.Service)
		return currentService, nil
	}
	return nil, nil
}