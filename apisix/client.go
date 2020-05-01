package apisix

import (
	"fmt"
	"net/http"
	"time"
	
	"gopkg.in/resty.v1"
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
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", resp.StatusCode(), resp.Body())
	}
	return resp.Body(), nil
}