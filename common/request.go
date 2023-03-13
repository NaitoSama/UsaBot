package common

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// RequestTo 发送请求 记得关闭resp的body
func RequestTo(url string, method string, contentType string, payload interface{}) (*http.Response, error) {
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
		req, err = http.NewRequest(method, url, body)
		if err != nil {
			return nil, err
		}
	} else if contentType == "" {
		req, err = http.NewRequest(method, url, nil)
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
