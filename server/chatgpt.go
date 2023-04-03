package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"UsaBot/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
		proxyUrl, _ := url.Parse(config.Config.General.Proxy)
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
		common.Logln(2, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Config.ChatGPT.AccessToken)
	res, err := client.Do(req)
	if err != nil {
		common.Logln(2, err)
		return
	}
	defer res.Body.Close()
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		common.Logln(2, err)
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}
	chatGPTresponse := Models.ChatGPTResponse{}
	err = json.Unmarshal(resData, &chatGPTresponse)
	if err != nil {
		common.Logln(2, err)
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

func ChatWithContext(msg Models.Message, user Models.ChatGPTUserInfo) {
	var role string
	var client *http.Client
	var messages []Models.ChatGPTMessage

	if strings.Contains(msg.Message, "&#91;设定&#93;") {
		role = "system"
	} else {
		role = "user"
	}
	var contexts []Models.ChatGPTContext
	Models.DB.Model(&Models.ChatGPTContext{}).Where("user = ? and state = ?", user.User, "enable").Order("id desc").Find(&contexts)
	//if len(contexts) >= user.MaxContexts {
	//	common.ErrorResponse(true, msg.GroupID, errors.New("您的上下文数量已达上限,现在为您进行清除"))
	//	Models.DB.Model(&Models.ChatGPTContext{}).Where("user = ? and state = ?", user.User, "enable").Update("state", "disable")
	//}

	message := msg.Message[strings.Index(msg.Message, "]")+1:]
	msgContent := Models.ChatGPTMessage{
		Role:    role,
		Content: message,
	}

	for _, v := range contexts {
		messages = append(messages, Models.ChatGPTMessage{
			Role:    v.Role,
			Content: v.Content,
		})
	}

	messages = append(messages, msgContent)

	content := Models.ChatGPT{
		Model:    config.Config.ChatGPT.Model,
		Messages: messages,
	}
	data, _ := json.Marshal(content)
	param := bytes.NewBuffer(data)

	if config.Config.ChatGPT.UseProxy {
		proxyUrl, _ := url.Parse(config.Config.General.Proxy)
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
		common.Logln(2, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Config.ChatGPT.AccessToken)
	res, err := client.Do(req)
	if err != nil {
		common.Logln(2, err)
		return
	}
	defer res.Body.Close()
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		common.Logln(2, err)
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}
	chatGPTresponse := Models.ChatGPTResponse{}
	err = json.Unmarshal(resData, &chatGPTresponse)
	if err != nil {
		common.Logln(2, err)
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}

	record := Models.ChatGPTContext{
		Role:    role,
		Content: messages[len(messages)-1].Content,
		User:    msg.Sender.UserID,
		State:   "enable",
	}
	Models.DB.Create(&record)

	for _, v := range chatGPTresponse.Choices {
		temp := fmt.Sprintf("[CQ:at,qq=%d] %s", msg.Sender.UserID, v.Message.Content)
		replyContent := Models.SendGroupMessage{
			GroupID: msg.GroupID,
			Message: temp,
		}
		common.PostToCQHTTPNoResponse(replyContent, "/send_group_msg")
		record = Models.ChatGPTContext{
			Role:    "assistant",
			Content: v.Message.Content,
			User:    msg.Sender.UserID,
			State:   "enable",
		}
		Models.DB.Create(&record)
	}
}
