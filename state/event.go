package state

import (
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
)

type ApisixCombination struct {
	Routes    []*v1.Route
	Services  []*v1.Service
	Upstreams []*v1.Upstream
}

type RouteCompare struct {
	OldRoutes []*v1.Route
	NewRoutes []*v1.Route
}

type Quit struct {
	Name string
}

const (
	RouteKind    = "route"
	ServiceKind  = "service"
	UpstreamKind = "upstream"
	Create       = "create"
	Update       = "update"
	Delete       = "delete"
)

type Event struct {
	Kind string      // route/service/upstream
	Op   string      // create update delete
	Obj  interface{} // the obj of kind
}
