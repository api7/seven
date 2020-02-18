package state

import (
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
	"github.com/gxthrj/seven/apisix"
	"github.com/gxthrj/seven/DB"
)

var UpstreamQueue chan UpstreamQueueObj
var ServiceQueue chan ServiceQueueObj

func init(){
	UpstreamQueue =  make(chan UpstreamQueueObj, 500)
	ServiceQueue =  make(chan ServiceQueueObj, 500)
	go WatchUpstream()
	go WatchService()
}

func WatchService(){
	for{
		sqo := <- ServiceQueue
		// solver service
		SolverService(sqo.Services, sqo.RouteWorkerGroup)
	}
}

func WatchUpstream(){
	for{
		uqo := <- UpstreamQueue
		SolverUpstream(uqo.Upstreams, uqo.ServiceWorkerGroup)
	}
}

// Solver
func (s *ApisixCombination)Solver() (bool, error){
	// 1.route workers
	rwg := NewRouteWorkers(s.Routes)
	// 2.service workers
	swg := NewServiceWorkers(s.Services, &rwg)
	//sqo := &ServiceQueueObj{Services: s.Services, RouteWorkerGroup: rwg}
	//sqo.AddQueue()
	// 3.upstream workers
	uqo := &UpstreamQueueObj{Upstreams: s.Upstreams, ServiceWorkerGroup: swg}
	uqo.AddQueue()
	return true, nil
}
// UpstreamQueueObj for upstream queue
type UpstreamQueueObj struct {
	Upstreams []*v1.Upstream
	ServiceWorkerGroup ServiceWorkerGroup
}

// AddQueue make upstreams in order
// upstreams is group by CRD
func (uqo *UpstreamQueueObj) AddQueue(){
	UpstreamQueue <- *uqo
}

type ServiceQueueObj struct {
	Services []*v1.Service
	RouteWorkerGroup RouteWorkerGroup
}

// AddQueue make upstreams in order
// upstreams is group by CRD
func (sqo *ServiceQueueObj) AddQueue(){
	ServiceQueue <- *sqo
}

// Sync remove from apisix
func (rc *RouteCompare) Sync() error{
	for _, old := range rc.OldRoutes{
		needToDel := true
		for _, nr := range rc.NewRoutes {
			if old.Name == nr.Name {
				needToDel = false
				break
			}
		}
		if needToDel {
			request := DB.RouteRequest{Name: *old.Name}

			if route, err := request.FindByName(); err != nil {
				// log error
			}else {
				if err = apisix.DeleteRoute(route); err == nil {
					db := DB.RouteDB{Routes: []*v1.Route{route}}
					db.DeleteRoute()
				}

			}

		}
	}
	return nil
}