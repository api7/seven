package state

import "github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"

type routeWorker struct {
	*v1.Route
	Event chan Event
	Quit  chan Quit
}

// RouteWorkerGroup for broadcast from service to route
type RouteWorkerGroup map[string][]*routeWorker

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
