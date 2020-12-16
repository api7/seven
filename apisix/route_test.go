package apisix

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteUnmarshalJSON(t *testing.T) {
	var route Routes
	emptyData := `
{
	"key": "test",
	"nodes": {}
}
`
	err := json.Unmarshal([]byte(emptyData), &route)
	assert.Nil(t, err)

	emptyData = `
{
	"key": "test",
	"nodes": {"a": "b", "c": "d"}
}
`
	err = json.Unmarshal([]byte(emptyData), &route)
	assert.Equal(t, err.Error(), "unexpected non-empty object")

	emptyArray := `
{
	"key": "test",
	"nodes": []
}
`
	err = json.Unmarshal([]byte(emptyArray), &route)
	assert.Nil(t, err)

	normalData := `
{
	"key": "test",
	"nodes": [
		{
			"key": "route 1",
			"value": {
				"desc": "test route 1",
				"upstream_id": "123",
				"service_id": "12345",
				"host": "foo.com",
				"uri": "/bar/baz",
				"methods": ["GET", "POST"]
			}
		}
	]
}
`
	err = json.Unmarshal([]byte(normalData), &route)
	assert.Nil(t, err)
	assert.Equal(t, route.Key, "test")
	assert.Equal(t, len(route.Routes), 1)

	key := *route.Routes[0].Key
	assert.Equal(t, key, "route 1")
	desc := *route.Routes[0].Value.Desc
	assert.Equal(t, desc, "test route 1")
	upstreamId := *route.Routes[0].Value.UpstreamId
	assert.Equal(t, upstreamId, "123")
	svcId := *route.Routes[0].Value.ServiceId
	assert.Equal(t, svcId, "12345")
	assert.Equal(t, *route.Routes[0].Value.Host, "foo.com")
	assert.Equal(t, *route.Routes[0].Value.Uri, "/bar/baz")
	assert.Equal(t, *route.Routes[0].Value.Methods[0], "GET")
	assert.Equal(t, *route.Routes[0].Value.Methods[1], "POST")
}
