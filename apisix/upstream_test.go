package apisix

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpstreamsUnmarshalJSON(t *testing.T) {
	var ups Upstreams
	emptyData := `
{
	"key": "test",
	"nodes": {}
}
`
	err := json.Unmarshal([]byte(emptyData), &ups)
	assert.Nil(t, err)

	emptyData = `
{
	"key": "test",
	"nodes": {"a": "b", "c": "d"}
}
`
	err = json.Unmarshal([]byte(emptyData), &ups)
	assert.Equal(t, err.Error(), "unexpected non-empty object")

	emptyArray := `
{
	"key": "test",
	"nodes": []
}
`
	err = json.Unmarshal([]byte(emptyArray), &ups)
	assert.Nil(t, err)

	normalData := `
{
	"key": "test",
	"nodes": [
		{
			"key": "ups1",
			"value": {
				"desc": "test upstream 1",
				"type": "rr",
				"nodes": {
					"192.168.12.12": 100
				}
			}
		}
	]
}
`
	err = json.Unmarshal([]byte(normalData), &ups)
	assert.Nil(t, err)
	assert.Equal(t, ups.Key, "test")
	assert.Equal(t, len(ups.Upstreams), 1)

	key := *ups.Upstreams[0].Key
	assert.Equal(t, key, "ups1")
	desc := *ups.Upstreams[0].UpstreamNodes.Desc
	assert.Equal(t, desc, "test upstream 1")
	lb := *ups.Upstreams[0].UpstreamNodes.LBType
	assert.Equal(t, lb, "rr")

	assert.Equal(t, len(ups.Upstreams[0].UpstreamNodes.Nodes), 1)
	assert.Equal(t, ups.Upstreams[0].UpstreamNodes.Nodes["192.168.12.12"], int64(100))
}
