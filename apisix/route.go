package apisix

import (
	"encoding/json"
	"fmt"
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
	"github.com/gxthrj/seven/utils"
	"strings"
	"github.com/gxthrj/seven/conf"
	"github.com/gxthrj/seven/DB"
)

// FindCurrentRoute find current route in memDB
func FindCurrentRoute(route *v1.Route) (*v1.Route,error){
	db := &DB.RouteRequest{Name: *(route.Name)}
	currentRoute, _ := db.FindByName()
	if currentRoute != nil {
		return currentRoute, nil
	} else {
		// find from apisix
		if routes, err := ListRoute(); err != nil {
			// todo log error
		} else {
			for _, r := range routes {
				if r.Name !=nil && *r.Name == *route.Name {
					// insert to memDB
					db := &DB.RouteDB{Routes: []*v1.Route{r}}
					db.Insert()
					// return
					return r, nil
				}
			}
		}

	}
	return nil, fmt.Errorf("NOT FOUND")
}

// ListRoute list route from etcd , convert to v1.Route
func ListRoute() ([]*v1.Route, error) {
	url := conf.BaseUrl + "/routes"
	ret, _ := Get(url)
	var routesResponse RoutesResponse
	if err := json.Unmarshal(ret, &routesResponse); err != nil {
		return nil, fmt.Errorf("json转换失败")
	} else {
		routes := make([]*v1.Route, 0)
		for _, u := range routesResponse.Routes.Routes {
			if n, err := u.convert(); err == nil {
				routes = append(routes, n)
			} else {
				return nil, fmt.Errorf("upstream: %s 转换失败, %s", *u.Value.Desc, err.Error())
			}
		}
		return routes, nil
	}
}

func AddRoute(route *v1.Route, baseUrl string) (*RouteResponse, error) {
	url := fmt.Sprintf("%s/routes", baseUrl)
	rr := convert2RouteRequest(route)
	if b, err := json.Marshal(rr); err != nil {
		return nil, err
	} else {
		if res, err := utils.Post(url, b); err != nil {
			return nil, err
		} else {
			var routeResp RouteResponse
			if err = json.Unmarshal(res, &routeResp); err != nil {
				return nil, err
			} else {
				return &routeResp, nil
			}
		}
	}
}

func UpdateRoute(route *v1.Route, baseUrl string) error {
	url := fmt.Sprintf("%s/routes/%s", baseUrl, *route.ID)
	rr := convert2RouteRequest(route)
	if b, err := json.Marshal(rr); err != nil {
		return err
	} else {
		if _, err := utils.Patch(url, b); err != nil {
			return err
		} else {
			return nil
		}
	}
}

type Redirect struct {
	RetCode int64 `json:"ret_code"`
	Uri string `json:"uri"`
}

func convert2RouteRequest(route *v1.Route) *RouteRequest {
	return &RouteRequest{
		Desc:      *route.Name,
		Host:      *route.Host,
		Uri:       *route.Path,
		ServiceId: *route.ServiceId,
		Plugins:   route.Plugins,
	}
}

// convert apisix RouteResponse -> apisix-types v1.Route
func (r *Route) convert() (*v1.Route, error) {
	// id
	key := r.Key
	ks := strings.Split(*key, "/")
	id := ks[len(ks)-1]
	// name
	name := r.Value.Desc
	// host
	host := r.Value.Host
	// path
	path := r.Value.Uri
	// method
	methods := r.Value.Methods
	// upstreamId
	upstreamId := r.Value.UpstreamId
	// serviceId
	serviceId := r.Value.ServiceId
	// plugins
	var plugins v1.Plugins
	plugins = r.Value.Plugins

	return &v1.Route{
		ID:         &id,
		Name:       name,
		Host:       host,
		Path:       path,
		Methods:    methods,
		UpstreamId: upstreamId,
		ServiceId:  serviceId,
		Plugins:    &plugins,
	}, nil
}

type RoutesResponse struct {
	Routes Routes `json:"node"`
}

type Routes struct {
	Key    string  `json:"key"`
	Routes []Route `json:"nodes"`
}

type RouteResponse struct {
	Action string `json:"action"`
	Route  Route  `json:"node"`
}

type Route struct {
	Key   *string `json:"key"`   // route key
	Value Value  `json:"value"` // route content
}

type Value struct {
	UpstreamId *string                `json:"upstream_id"`
	ServiceId  *string                `json:"service_id"`
	Plugins    map[string]interface{} `json:"plugins"`
	Host       *string                `json:"host,omitempty"`
	Uri        *string                `json:"uri"`
	Desc       *string                `json:"desc"`
	Methods    []*string              `json:"methods,omitempty"`
}

type RouteRequest struct {
	Desc      string      `json:"desc,omitempty"`
	Uri       string      `json:"uri,omitempty"`
	Host      string      `json:"host,omitempty"`
	ServiceId string      `json:"service_id,omitempty"`
	Plugins   *v1.Plugins `json:"plugins,omitempty"`
}
