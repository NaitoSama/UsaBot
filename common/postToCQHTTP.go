package common

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

const url = "http://127.0.0.1:5700"
const accessToken = ""
const contentType = "application/json"
const authorization = ""

var lock sync.RWMutex

func PostToCQHTTPNoResponse(content interface{}, path string) {
	configData, _ := json.Marshal(content)
	param := bytes.NewBuffer(configData)
	client := http.DefaultClient
	lock.RLock()
	req, err := http.NewRequest("POST", url+path, param)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", authorization)
	lock.RUnlock()
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
}

func PostToCQHTTPWithResponse(content interface{}, path string) (*http.Response, error) {
	configData, _ := json.Marshal(content)
	param := bytes.NewBuffer(configData)
	client := http.DefaultClient
	lock.RLock()
	req, err := http.NewRequest("POST", url+path, param)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", authorization)
	lock.RUnlock()
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return res, nil
}
