package state

import (
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
	"github.com/gxthrj/seven/apisix"
	"strconv"
	"github.com/gxthrj/seven/utils"
)

var BaseUrl = "http://172.16.20.90:30116/apisix/admin"

func conf(baseUrl string){
	BaseUrl = baseUrl
}

// ListFromApisix list all object from apisix
func ListFromApisix(){

}

// InitDB insert object into memDB first time
func InitDB(){
	routes, _ := apisix.ListRoute(BaseUrl)
	upstreams, _ := apisix.ListUpstream(BaseUrl)
	apisix.InsertRoute(routes)
	apisix.InsertUpstreams(upstreams)
}

// LoadTargetState load targetState from ... maybe k8s CRD
func LoadTargetState(routes []*v1.Route, upstreams []*v1.Upstream){

	// 1.diff
	// 2.send event
}

// paddingRoute padding route from memDB
func paddingRoute(route *v1.Route, currentRoute *v1.Route){
	// 1.get object from memDB
	// 2.padding object, just id
	if currentRoute == nil {
		// NOT FOUND : set Id = 0
		id := strconv.Itoa(0)
		route.ID = &id
	} else {
		route.ID = currentRoute.ID
	}
}

// paddingUpstream padding upstream from memDB
func paddingUpstream(upstream *v1.Upstream, currentUpstream *v1.Upstream){
	if currentUpstream == nil {
		// NOT FOUND : set Id = 0
		id := strconv.Itoa(0)
		upstream.ID = &id
	} else {
		upstream.ID = currentUpstream.ID
	}
}

// NewRouteWorkers make routeWrokers group by service per CRD
// 1.make routes group by (1_2_3) it may be a map like map[1_2_3][]Route;
// 2.route is listenning Event from the ready of 1_2_3;
func NewRouteWorkers(routes []*v1.Route, quit chan Quit) RouteWorkerGroup{
	rwg := make(RouteWorkerGroup)
	for _, r := range routes {
		rw := &routeWorker{Route: r, Quit: quit}
		rw.start()
		rwg.Add(*r.ServiceName, rw)
	}
	return rwg
}

// 3.route get the Event and trigger a padding for object,then diff,sync;
func (r *routeWorker) trigger(event Event) error{
	// todo consumer Event

	// padding
	currentRoute, _ := apisix.FindRoute(r.Route)
	paddingRoute(r.Route, currentRoute)
	// diff
	hasDiff, err := utils.HasDiff(r.Route, currentRoute)
	// sync
	if err != nil {
		return err
	}
	if hasDiff {
		r.sync()
	}
	// todo broadcast

	return nil
}

// sync
func (r *routeWorker) sync(){
	if *r.Route.ID != strconv.Itoa(0) {
		apisix.UpdateRoute(r.Route, BaseUrl)
	} else {
		apisix.AddRoute(r.Route, BaseUrl)
	}
}

// service
func NewServiceWorkers(services []*v1.Service, quit chan Quit) ServiceWorkerGroup{
	swg := make(ServiceWorkerGroup)
	for _, s := range services {
		rw := &serviceWorker{Service: s, Quit: quit}
		rw.start()
		swg.Add(*s.UpstreamId, rw) // todo key is upstreamName
	}
	return swg
}

// upstream
func SolverUpstream(upstreams []*v1.Upstream){
	for _, u := range upstreams {
		if currentUpstream, err := apisix.FindUpstreamByName(*u.Name); err != nil {
			return
		} else {
			paddingUpstream(u, currentUpstream)
			// diff
			hasDiff, _ := utils.HasDiff(u, currentUpstream)
			if hasDiff {
				// sync
				if *u.ID != strconv.Itoa(0) {
					apisix.UpdateUpstream(u, BaseUrl)
				} else {
					apisix.AddUpstream(u, BaseUrl)
				}
			}
		}
		// todo broadcast
	}
}

