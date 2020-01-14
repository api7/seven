package apisix

import (
	"github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
	"encoding/json"
	"fmt"
	"strings"
	"strconv"
)



// List list upstream from etcd , convert to v1.Upstream
func List(baseUrl string) ([]*v1.Upstream, error) {
	url := baseUrl + "/upstreams"
	ret, _ := Get(url)
	var upstreamsResponse UpstreamsResponse
	if err := json.Unmarshal(ret, &upstreamsResponse); err != nil {
		return nil, fmt.Errorf("json转换失败")
	} else {
		upstreams := make([]*v1.Upstream, len(upstreamsResponse.Upstreams.Upstreams))
		for _, u := range upstreamsResponse.Upstreams.Upstreams {
			if n, err := u.convert(); err == nil {
				upstreams = append(upstreams, n)
			} else {
				return nil, fmt.Errorf("upstream: %s 转换失败, %s", u.UpstreamNodes.Desc, err.Error())
			}
		}
		return upstreams, nil
	}
}

// convert convert Upstream from etcd to v1.Upstream
func (u *Upstream)convert() (*v1.Upstream, error){
	// id
	keys := strings.Split(u.Key, "/")
	id := keys[len(keys) - 1]
	// Name
	name := u.UpstreamNodes.Desc
	// type
	LBType := u.UpstreamNodes.LBType
	// key
	key := u.Key
	// nodes
	nodes := make([]*v1.Node, len(u.UpstreamNodes.Nodes))
	for k, v := range u.UpstreamNodes.Nodes {
		ks := strings.Split(k, ":")
		ip := ks[0]
		port, _ := strconv.Atoi(ks[1])
		weight := int(v)
		node := &v1.Node{IP: &ip, Port: &port, Weight: &weight}
		nodes = append(nodes, node)
	}

	return &v1.Upstream{ID: &id, Name: &name, Type: &LBType, Key: &key, Nodes: nodes}, nil
}

type UpstreamsResponse struct {
	Upstreams Upstreams `json:"node"`
}

type Upstreams struct{
	Key string `json:"key"` // 用来定位upstreams 列表
	Upstreams []Upstream `json:"nodes"`
}

type Upstream struct {
	Key string `json:"key"` // upstream key
	UpstreamNodes UpstreamNodes `json:"value"`
}

type UpstreamNodes struct {
	Nodes map[string]int64 `json:"nodes"`
	Desc string `json:"desc"` // upstream name  = k8s svc
	LBType string `json:"type"` // 负载均衡类型
}