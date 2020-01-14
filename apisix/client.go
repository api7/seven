package apisix

import (
	"gopkg.in/resty.v1"
	"time"
)
const (
	timeout = 3000
)

func Get(url string) ([]byte, error){
	r := resty.New().
		SetTimeout(time.Duration(timeout)*time.Millisecond).
		R().
		SetHeader("content-type", "application/json")
	resp, err := r.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}