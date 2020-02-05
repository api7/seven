package apisix

import (
	"encoding/json"
	"fmt"
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
	"github.com/gxthrj/seven/utils"
	"github.com/golang/glog"
)

func AddService(service *v1.Service, baseUrl string) (*ServiceResponse, error) {
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
				return &uRes, nil
			}
		}
	}
}

func UpdateService(service *v1.Service, baseUrl string) (*ServiceResponse, error) {
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
				return &uRes, nil
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
