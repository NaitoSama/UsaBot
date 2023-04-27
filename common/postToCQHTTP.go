package common

import (
	"UsaBot/Models"
	"UsaBot/config"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
)

var cqurl = config.Config.General.CQHttpUrl

const accessToken = ""
const contentType = "application/json"
const authorization = ""

var lock sync.RWMutex

func PostToCQHTTPNoResponse(content interface{}, path string) {
	configData, _ := json.Marshal(content)
	param := bytes.NewBuffer(configData)
	client := http.DefaultClient
	lock.RLock()
	req, err := http.NewRequest("POST", cqurl+path, param)
	if err != nil {
		Logln(2, err)
		return
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", authorization)
	lock.RUnlock()
	res, err := client.Do(req)
	if err != nil {
		Logln(2, err)
		return
	}
	defer res.Body.Close()
}

func PostToCQHTTPWithResponse(content interface{}, path string) (*http.Response, error) {
	configData, _ := json.Marshal(content)
	param := bytes.NewBuffer(configData)
	client := http.DefaultClient
	lock.RLock()
	req, err := http.NewRequest("POST", cqurl+path, param)
	if err != nil {
		Logln(2, err)
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", authorization)
	lock.RUnlock()
	res, err := client.Do(req)
	if err != nil {
		Logln(2, err)
		return nil, err
	}
	return res, nil
}

// GroupChatSender 发送群聊，返回是否成功发送与错误信息
func GroupChatSender(groupID int64, content string) (bool, error) {
	message := Models.SendGroupMessage{
		GroupID: groupID,
		Message: content,
	}
	response, err := PostToCQHTTPWithResponse(message, "/send_group_msg")
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	respData, err := io.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	respStruct := Models.SendGroupMessageResponse{}
	err = json.Unmarshal(respData, &respStruct)

	if respStruct.Status != "ok" {
		return false, errors.New(respStruct.Wording)
	} else {
		return true, nil
	}
}
