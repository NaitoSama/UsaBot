package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"UsaBot/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func ChatGPT(msg Models.Message, role string) {
	var client *http.Client

	message := msg.Message[strings.Index(msg.Message, "]")+1:]
	msgContent := Models.ChatGPTMessage{
		Role:    role,
		Content: message,
	}
	content := Models.ChatGPT{
		Model:    config.Config.ChatGPT.Model,
		Messages: []Models.ChatGPTMessage{msgContent},
	}
	data, _ := json.Marshal(content)
	param := bytes.NewBuffer(data)

	if config.Config.ChatGPT.UseProxy {
		proxyUrl, _ := url.Parse(config.Config.ChatGPT.Proxy)
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		client = &http.Client{
			Transport: transport,
		}
	} else {
		client = &http.Client{}
	}

	req, err := http.NewRequest("POST", config.Config.ChatGPT.Url, param)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Config.ChatGPT.AccessToken)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}
	chatGPTresponse := Models.ChatGPTResponse{}
	err = json.Unmarshal(resData, &chatGPTresponse)
	if err != nil {
		log.Println(err)
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}

	for _, v := range chatGPTresponse.Choices {
		temp := fmt.Sprintf("[CQ:at,qq=%d] %s", msg.Sender.UserID, v.Message.Content)
		replyContent := Models.SendGroupMessage{
			GroupID: msg.GroupID,
			Message: temp,
		}
		common.PostToCQHTTPNoResponse(replyContent, "/send_group_msg")
	}
}
