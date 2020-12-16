package apisix

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSslUnmarshalJSON(t *testing.T) {
	var sslList SslList
	emptyData := `
{
	"key": "test",
	"nodes": {}
}
`
	err := json.Unmarshal([]byte(emptyData), &sslList)
	assert.Nil(t, err)

	notEmptyObject := `
{
	"key": "test",
	"nodes": {"a": "b", "c": "d"}
}
`
	err = json.Unmarshal([]byte(notEmptyObject), &sslList)
	assert.Equal(t, err.Error(), "unexpected non-empty object")

	emptyArray := `
{
	"key": "test",
	"nodes": []
}
`
	err = json.Unmarshal([]byte(emptyArray), &sslList)
	assert.Nil(t, err)

	normalData := `
{
	"key": "test",
	"nodes": [
		{
			"key": "ssl id",
			"value": {
				"snis": ["test.apisix.org"],
				"cert": "root",
				"key": "123456",
				"status": 1
			}
		}
	]
}
`
	err = json.Unmarshal([]byte(normalData), &sslList)
	assert.Nil(t, err)
	assert.Equal(t, len(sslList.SslNodes), 1)

	key := *sslList.SslNodes[0].Key
	assert.Equal(t, key, "ssl id")
	cert := *sslList.SslNodes[0].Ssl.Cert
	assert.Equal(t, cert, "root")
	sslKey := *sslList.SslNodes[0].Ssl.Key
	assert.Equal(t, sslKey, "123456")
	sni := *sslList.SslNodes[0].Ssl.Snis[0]
	assert.Equal(t, sni, "test.apisix.org")
}
