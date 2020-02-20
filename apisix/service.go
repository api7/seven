package apisix

import (
	"encoding/json"
	"fmt"
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
	"github.com/gxthrj/seven/utils"
	"github.com/golang/glog"
	"strings"
	"github.com/gxthrj/seven/DB"
	"github.com/gxthrj/seven/conf"
)

// FindCurrentService find service from memDB,
// if Not Found, find service from apisix
func FindCurrentService(group, name, fullName string) (*v1.Service, error){
	db := DB.ServiceRequest{Group: group, Name: name, FullName: fullName}
	currentService, _ := db.FindByName()
	if currentService != nil {
		return currentService, nil
	}else {
		// find service from apisix
		if services, err := ListService(group); err != nil {
			// todo log error
			glog.Info(err.Error())
		}else {
			for _, s := range services {
				if s.Name != nil && *(s.Name) == name {
					// and save to memDB
					db := &DB.ServiceDB{Services: []*v1.Service{s}}
					db.Insert()
					// return
					return s, nil
				}
			}
		}
	}
	return nil, nil
}

// ListUpstream list upstream from etcd , convert to v1.Upstream
func ListService (group string) ([]*v1.Service, error) {
	baseUrl := conf.FindUrl(group)
	url := baseUrl + "/services"
	ret, _ := Get(url)
	var servicesResponse ServicesResponse
	if err := json.Unmarshal(ret, &servicesResponse); err != nil {
		return nil, fmt.Errorf("json转换失败")
	} else {
		result := make([]*v1.Service, 0)
		for _, u := range servicesResponse.Services.Services {
			if n, err := u.convert(group); err == nil {
				result = append(result, n)
			} else {
				return nil, fmt.Errorf("service : %s 转换失败, %s", *u.ServiceValue.Desc, err.Error())
			}
		}
		return result, nil
	}
}

// convert convert Service from etcd to v1.Service
func (u *Service)convert(group string) (*v1.Service, error){
	// id
	keys := strings.Split(*u.Key, "/")
	id := keys[len(keys) - 1]
	// Name
	name := u.ServiceValue.Desc
	// upstreamId
	upstreamId := u.ServiceValue.UpstreamId
	// plugins
	plugins := &v1.Plugins{}
	for k, v := range u.ServiceValue.Plugins {
		(*plugins)[k] = v
	}
	fullName := *name
	if group != ""{
		fullName = group + "_" + *name
	}
	return &v1.Service{ID: &id, FullName: &fullName, Group: &group, Name: name, UpstreamId: upstreamId, Plugins: plugins}, nil
}

func AddService(service *v1.Service) (*ServiceResponse, error) {
	baseUrl := conf.FindUrl(*service.Group)
	url := fmt.Sprintf("%s/services", baseUrl)
	ur := convert2ServiceRequest(service)
	if b, err := json.Marshal(ur); err != nil {
		return nil, err
	} else {
		if res, err := utils.Post(url, b); err != nil {
			return nil, err
		} else {
			var uRes ServiceResponse
			if err = json.Unmarshal(res, &uRes); err != nil {
				return nil, err
			} else {
				if uRes.Service.Key != nil {
					return &uRes, nil
				} else {
					return nil, fmt.Errorf("apisix service not expected response")
				}

			}
		}
	}
}

func UpdateService(service *v1.Service) (*ServiceResponse, error) {
	baseUrl := conf.FindUrl(*service.Group)
	url := fmt.Sprintf("%s/services/%s", baseUrl, *service.ID)
	ur := convert2ServiceRequest(service)
	if b, err := json.Marshal(ur); err != nil {
		return nil, err
	} else {
		if res, err := utils.Patch(url, b); err != nil {
			return nil, err
		} else {
			var uRes ServiceResponse
			if err = json.Unmarshal(res, &uRes); err != nil {
				return nil, err
			} else {
				if uRes.Service.Key != nil {
					return &uRes, nil
				} else {
					var errResp ErrorResponse
					json.Unmarshal(res, &errResp)
					glog.Error(errResp.Message)
					return nil, fmt.Errorf("apisix service not expected response %s", errResp.Message)
				}
			}
		}
	}
}

func convert2ServiceRequest(service *v1.Service) *ServiceRequest {
	request := &ServiceRequest{
		Desc:       service.Name,
		UpstreamId: service.UpstreamId,
		Plugins:    service.Plugins,
	}
	glog.Info(*request.Desc)
	return request
}

type ServiceRequest struct {
	Desc       *string                `json:"desc,omitempty"`
	UpstreamId *string                `json:"upstream_id"`
	Plugins    *v1.Plugins `json:"plugins,omitempty"`
}


type ServicesResponse struct {
	Services Services `json:"node"`
}

type Services struct{
	Key string `json:"key"` // 用来定位upstreams 列表
	Services []Service `json:"nodes"`
}

type ServiceResponse struct {
	Action  string  `json:"action"`
	Service Service `json:"node"`
}

type Service struct {
	Key          *string      `json:"key"` // service key
	ServiceValue ServiceValue `json:"value,omitempty"`
}

type ServiceValue struct {
	UpstreamId *string                `json:"upstream_id,omitempty"`
	Plugins    map[string]interface{} `json:"plugins"`
	Desc       *string                `json:"desc,omitempty"`
}
