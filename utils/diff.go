package utils

import (
	"github.com/yudai/gojsondiff"
	"encoding/json"
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
		return d.Modified(), nil
	}
}
