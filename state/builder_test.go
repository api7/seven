package state

import (
	"testing"
	"github.com/gxthrj/seven/utils"
)

type school struct {
	*province
	Name string `json:"name"`
	Address string `json:"address"`
}

type province struct {
	Location string `json:"location"`
}

func Test_diff(t *testing.T){
	//p1 := &province{Location: "jiangsu"}
	p2 := &province{Location: "zh"}
	s1 := &school{Name: "hello", Address: "this is a address"}
	s2 := &school{Name: "hello", Address: "this is a address", province: p2}
	t.Log(s1)
	t.Log(s2)
	if d, err := utils.Diff(s1, s2); err != nil {
		t.Log(err.Error())
	}else {
		//t.Logf("s1 vs s2 hasDiff ? %v", d)
		t.Log(d)
		for _, delta := range d.Deltas() {
			t.Log(delta.Similarity())
		}

	}

}
