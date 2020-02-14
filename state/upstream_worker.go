package state

import "github.com/gxthrj/apisix-types/pkg/apis/apisix/v1"

type upstreamWorker struct {
	*v1.Upstream
	Event chan Event
	Quit  chan Quit
}

