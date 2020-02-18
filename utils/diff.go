package utils

import (
	"github.com/yudai/gojsondiff"
	"encoding/json"
	"github.com/golang/glog"
)

var (
	differ = gojsondiff.New()
)

func HasDiff(a, b interface{}) (bool, error){
	aJSON, err := json.Marshal(a)
	if err != nil {
		return false, err
	}
	bJSON, err := json.Marshal(b)
	if err != nil {
		return false, err
	}
	if d, err := differ.Compare(aJSON, bJSON); err != nil {
		return false, err
	}else {
		glog.Info(d.Deltas())
		return d.Modified(), nil
	}
}

func Diff(a, b interface{}) (gojsondiff.Diff, error){
	aJSON, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	bJSON, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	if d, err := differ.Compare(aJSON, bJSON); err != nil {
		return nil, err
	}else {
		return d, nil
	}
}
