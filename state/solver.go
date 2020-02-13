package state

import "github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"

var UpstreamQueue chan UpstreamQueueObj
func init(){
	UpstreamQueue =  make(chan UpstreamQueueObj, 500)
	go WatchUpstream()
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
	// 3.upstream workers
	//SolverUpstream(s.Upstreams, swg)
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
