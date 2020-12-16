package utils

import (
	"fmt"
	"net/http"
	"time"

	"gopkg.in/resty.v1"
)

const timeout = 3000

func Post(url string, bytes []byte) ([]byte, error) {
	r := resty.New().
		SetTimeout(time.Duration(timeout)*time.Millisecond).
		R().
		SetHeader("content-type", "application/json")
	r.SetBody(bytes)
	resp, err := r.Post(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("status: %d, body: %s", resp.StatusCode(), resp.Body())
	}
	return resp.Body(), nil
}

func Put(url string, bytes []byte) ([]byte, error) {
	r := resty.New().
		SetTimeout(time.Duration(timeout)*time.Millisecond).
		R().
		SetHeader("content-type", "application/json")
	r.SetBody(bytes)
	resp, err := r.Put(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("status: %d, body: %s", resp.StatusCode(), resp.Body())
	}
	return resp.Body(), nil
}

func Patch(url string, bytes []byte) ([]byte, error) {
	r := resty.New().
		SetTimeout(time.Duration(timeout)*time.Millisecond).
		R().
		SetHeader("content-type", "application/json")
	r.SetBody(bytes)
	resp, err := r.Patch(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("status: %d, body: %s", resp.StatusCode(), resp.Body())
	}
	return resp.Body(), nil
}

func Delete(url string) ([]byte, error) {
	r := resty.New().
		SetTimeout(time.Duration(timeout)*time.Millisecond).
		R().
		SetHeader("content-type", "application/json")
	resp, err := r.Delete(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", resp.StatusCode(), resp.Body())
	}
	return resp.Body(), nil
}
