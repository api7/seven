package state

import (
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
	"github.com/gxthrj/seven/DB"
	"github.com/gxthrj/seven/apisix"
	"github.com/gxthrj/seven/conf"
	"github.com/golang/glog"
	"github.com/gxthrj/seven/utils"
	"strconv"
	"strings"
)

type serviceWorker struct {
	*v1.Service
	Event chan Event
	Quit  chan Quit
}

// ServiceWorkerGroup for broadcast from upstream to service
type ServiceWorkerGroup map[string][]*serviceWorker

// start start watch event
func (w *serviceWorker) start(rwg *RouteWorkerGroup) {
	w.Event = make(chan Event)
	go func() {
		for {
			select {
			case event := <-w.Event:
				w.trigger(event, rwg)
			case <-w.Quit:
				return
			}
		}
	}()
}

func (w *serviceWorker) trigger(event Event, rwg *RouteWorkerGroup) error {
	glog.Infof("1.service trigger from %s, %s", event.Op, event.Kind)
	defer close(w.Quit)
	// consumer Event set upstreamID
	upstream := event.Obj.(*v1.Upstream)
	glog.Infof("2.service trigger from %s, %s", event.Op, *upstream.Name)
	w.UpstreamId = upstream.ID

	op := Update
	// padding
	currentRoute, _ := apisix.FindCurrentService(*w.Service.Name)
	paddingService(w.Service, currentRoute)
	// diff
	hasDiff, err := utils.HasDiff(w.Service, currentRoute)
	// sync
	if err != nil {
		return err
	}
	if hasDiff {
		if *w.Service.ID == strconv.Itoa(0) {
			op = Create
			// 1. sync apisix and get id
			if serviceResponse, err := apisix.AddService(w.Service, conf.BaseUrl); err != nil {
				// todo log error
				glog.Info(err.Error())
			}else {
				tmp := strings.Split(*serviceResponse.Service.Key, "/")
				*w.Service.ID = tmp[len(tmp) - 1]
			}
			// 2. sync memDB
			db := &DB.ServiceDB{Services: []*v1.Service{w.Service}}
			db.Insert()
			glog.Infof("create service %s, %s", *w.Name, *w.UpstreamId)
		}else {
			op = Update
			// 1. sync memDB
			db := DB.ServiceDB{Services: []*v1.Service{w.Service}}
			if err := db.UpdateService(); err != nil {
				// todo log error
			}
			// 2. sync apisix
			apisix.UpdateService(w.Service, conf.BaseUrl)
			glog.Infof("update service %s, %s", *w.Name, *w.UpstreamId)
		}
	}
	// broadcast to route
	routeWorkers := (*rwg)[*w.Service.Name]
	for _, rw := range routeWorkers{
		event := &Event{Kind: ServiceKind, Op: op, Obj: w.Service}
		glog.Infof("send event %s, %s, %s", event.Kind, event.Op, *w.Service.Name)
		rw.Event <- *event
	}
	return nil
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