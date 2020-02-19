package conf

import "testing"

func Test_map(t *testing.T){
	m1 := make(map[string]string)
	m1["a"] = "aa"
	m1["b"] = "bb"
	t.Log(m1["c"] == "")
}
