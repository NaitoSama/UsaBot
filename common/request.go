package common

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// RequestTo 发送请求 记得关闭resp的body
func RequestTo(urlS string, method string, contentType string, payload interface{}) (*http.Response, error) {
	var (
		req *http.Request
		err error
	)

	if contentType == "application/json" {
		bodyData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body := bytes.NewReader(bodyData)
		req, err = http.NewRequest(method, urlS, body)
		if err != nil {
			return nil, err
		}
	} else if contentType == "" {
		req, err = http.NewRequest(method, urlS, nil)
		if err != nil {
			return nil, err
		}
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RequestTOWithProxy 通过代理发送请求 记得关闭resp的body
func RequestTOWithProxy(urlS string, method string, contentType string, payload interface{}) (*http.Response, error) {
	var (
		req *http.Request
		err error
	)

	if contentType == "application/json" {
		bodyData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body := bytes.NewReader(bodyData)
		req, err = http.NewRequest(method, urlS, body)
		if err != nil {
			return nil, err
		}
	} else if contentType == "" {
		req, err = http.NewRequest(method, urlS, nil)
		if err != nil {
			return nil, err
		}
	}

	proxyUrl, _ := url.Parse("http://localhost:1080")
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
