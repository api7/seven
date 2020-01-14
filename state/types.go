package state

import (
	apisixv1 "github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"
)

type ApisixRawState struct {
	Services  []*apisixv1.Service
	Routes    []*apisixv1.Route
	Upstreams []*apisixv1.Upstream
}
