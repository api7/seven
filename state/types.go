package state

import (
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
	"github.com/gxthrj/seven/apisix"
	"github.com/gxthrj/seven/utils"
)

type Quit struct {
	Name string
}

type Event struct {
	Kind string      // route/service/upstream
	Op   string      // create update delete
	Obj  interface{} // the obj of kind
}

type routeWorker struct {
	*v1.Route
	Event chan Event
	Quit  chan Quit
}

// start start watch event
func (w *routeWorker) start() {
	w.Event = make(chan Event)
	go func() {
		for {
			select {
			case event := <-w.Event:
				w.trigger(event)
			case <-w.Quit:
				return
			}
		}
	}()
}

type serviceWorker struct {
	*v1.Service
	Event chan Event
	Quit  chan Quit
}

// start start watch event
func (w *serviceWorker) start() {
	w.Event = make(chan Event)
	go func() {
		for {
			select {
			case event := <-w.Event:
				w.trigger(event)
			case <-w.Quit:
				return
			}
		}
	}()
}

func (w *serviceWorker) trigger(event Event) error{
	// todo consumer Event

	// padding
	currentRoute, _ := apisix.FindServiceByName(*w.Service.Name)
	//paddingService(w.Route, currentRoute)
	// diff
	hasDiff, err := utils.HasDiff(w.Service, currentRoute)
	// sync
	if err != nil {
		return err
	}
	if hasDiff {
		//w.sync()
	}
	// todo broadcast

	return nil
}

type upstreamWorker struct {
	*v1.Upstream
	Event chan Event
	Quit  chan Quit
}

// RouteWorkerGroup for broadcast from service to route
type RouteWorkerGroup map[string][]*routeWorker

// ServiceWorkerGroup for broadcast from upstream to service
type ServiceWorkerGroup map[string][]*serviceWorker

func (rg *RouteWorkerGroup) Add(key string, rw *routeWorker) {
	routes := (*rg)[key]
	if routes == nil {
		routes = make([]*routeWorker, 0)
	}
	routes = append(routes, rw)
	(*rg)[key] = routes
}

func (rg *RouteWorkerGroup) Delete(key string, route *routeWorker) {
	routes := (*rg)[key]
	result := make([]*routeWorker, 0)
	for _, r := range routes {
		if r.Name != route.Name {
			result = append(result, r)
		}
	}
	(*rg)[key] = result
}

func (swg *ServiceWorkerGroup) Add(key string, s *serviceWorker) {
	sws := (*swg)[key]
	if sws == nil {
		sws = make([]*serviceWorker, 0)
	}
	sws = append(sws, s)
	(*swg)[key] = sws
}

func (swg *ServiceWorkerGroup) Delete(key string, s *serviceWorker) {
	sws := (*swg)[key]
	result := make([]*serviceWorker, 0)
	for _, r := range sws {
		if r.Name != s.Name {
			result = append(result, r)
		}
	}
	(*swg)[key] = result
}

//type ServiceGroup map[string][]*v1.Service
//
//func (sg *ServiceGroup) Add(key string, service *v1.Service){
//	services := (*sg)[key]
//	if services == nil {
//		services = make([]*v1.Service, 0)
//	}
//	services = append(services, service)
//	(*sg)[key] = services
//}
//
//func (sg *ServiceGroup) Delete(key string, service *v1.Service){
//	services := (*sg)[key]
//	result := make([]*v1.Service, 0)
//	for _, s := range services{
//		if s.Name != service.Name{
//			result = append(result, s)
//		}
//	}
//	(*sg)[key] = result
//}
