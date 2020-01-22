package state

import (
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
	"github.com/gxthrj/seven/apisix"
	"strconv"
	"github.com/gxthrj/seven/utils"
	"strings"
	"github.com/gxthrj/seven/DB"
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
	// padding object, just id
	if currentRoute == nil {
		// NOT FOUND : set Id = 0
		id := strconv.Itoa(0)
		route.ID = &id
	} else {
		route.ID = currentRoute.ID
	}
}

// padding service from memDB
func paddingService(service *v1.Service, currentService *v1.Service){
	if currentService == nil {
		id := strconv.Itoa(0)
		service.ID = &id
	} else {
		service.ID = currentService.ID
	}
}

// paddingUpstream padding upstream from memDB
func paddingUpstream(upstream *v1.Upstream, currentUpstream *v1.Upstream){
	// padding id
	if currentUpstream == nil {
		// NOT FOUND : set Id = 0
		id := strconv.Itoa(0)
		upstream.ID = &id
	} else {
		upstream.ID = currentUpstream.ID
	}
	// todo padding nodes ? or sync nodes from crd ?
}

// NewRouteWorkers make routeWrokers group by service per CRD
// 1.make routes group by (1_2_3) it may be a map like map[1_2_3][]Route;
// 2.route is listenning Event from the ready of 1_2_3;
func NewRouteWorkers(routes []*v1.Route) RouteWorkerGroup{
	rwg := make(RouteWorkerGroup)
	for _, r := range routes {
		quit := make(chan Quit)
		rw := &routeWorker{Route: r, Quit: quit}
		rw.start()
		rwg.Add(*r.ServiceName, rw)
	}
	return rwg
}

// 3.route get the Event and trigger a padding for object,then diff,sync;
func (r *routeWorker) trigger(event Event) error{
	defer close(r.Quit)
	// consumer Event
	service := event.Obj.(v1.Service)
	r.ServiceId = service.ID

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
		// 1. sync memDB
		db := &DB.RouteDB{r.Route}
		if err := db.UpdateRoute(); err != nil {
			// todo log error
		}
		// 2. sync apisix
		apisix.UpdateRoute(r.Route, BaseUrl)
	} else {
		// 1. sync apisix and get id
		if res, err := apisix.AddRoute(r.Route, BaseUrl); err != nil {
			// todo log error
		} else {
			tmp := strings.Split(res.Routes.Key, "/")
			*r.ID = tmp[len(tmp) - 1]
		}
		// 2. sync memDB
		apisix.InsertRoute([]*v1.Route{r.Route})
	}
}

// service
func NewServiceWorkers(services []*v1.Service, quit chan Quit, rwg *RouteWorkerGroup) ServiceWorkerGroup{
	swg := make(ServiceWorkerGroup)
	for _, s := range services {
		rw := &serviceWorker{Service: s, Quit: quit}
		rw.start(rwg)
		swg.Add(*s.UpstreamName, rw)
	}
	return swg
}

// upstream
func SolverUpstream(upstreams []*v1.Upstream, swg ServiceWorkerGroup){
	for _, u := range upstreams {
		op := Update
		if currentUpstream, err := apisix.FindUpstreamByName(*u.Name); err != nil {
			// todo log error
		} else {
			paddingUpstream(u, currentUpstream)
			// diff
			hasDiff, _ := utils.HasDiff(u, currentUpstream)
			if hasDiff {
				if *u.ID != strconv.Itoa(0) {
					op = Update
					// 1.sync memDB
					upstreamDB := &DB.UpstreamDB{u}
					if err := upstreamDB.UpdateUpstream(); err != nil {
						// todo log error
					}
					// 2.sync apisix
					apisix.UpdateUpstream(u, BaseUrl)
				} else {
					op = Create
					// 1.sync apisix and get response
					if upstreamResponse, err := apisix.AddUpstream(u, BaseUrl); err != nil {
						// todo log error
					}else {
						tmp := strings.Split(*upstreamResponse.Upstream.Key, "/")
						*u.ID = tmp[len(tmp) - 1]
					}
					// 2.sync memDB
					apisix.InsertUpstreams([]*v1.Upstream{u})
				}
			}
		}
		// anyway, broadcast to service
		serviceWorkers := swg[*u.Name]
		for _, sw := range serviceWorkers{
			event := &Event{Kind: UpstreamKind, Op: op, Obj: u}
			sw.Event <- *event
		}
	}
}
