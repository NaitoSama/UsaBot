package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"UsaBot/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func ChatGPT(msg Models.Message, role string) {
	var client *http.Client
	lock.RLock()
	defer lock.RUnlock()

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
	lock.RLock()
	defer lock.RUnlock()
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

func chatGPTMainHandler(body Models.Message) {
	lock.Lock()
	user, ok := Models.ChatGPTUsers[body.Sender.UserID]
	if !ok {
		user = Models.ChatGPTUserInfo{
			User:          body.Sender.UserID,
			EnableContext: false,
			MaxContexts:   50,
		}
		Models.ChatGPTUsers[body.Sender.UserID] = user
		lock.Unlock()
		Models.DB.Create(&user)
	} else {
		lock.Unlock()
	}

	var count int64 = 0
	Models.DB.Model(Models.ChatGPTContext{}).Where("user = ? and state = ?", user.User, "enable").Count(&count)

	if count >= int64(user.MaxContexts) {
		common.ErrorResponse(true, body.GroupID, errors.New("您的上下文数量已达上限,现在为您进行清除"))
		Models.DB.Model(&Models.ChatGPTContext{}).Where("user = ? and state = ?", user.User, "enable").Update("state", "disable")
		count = 0
	}

	if strings.Contains(body.Message, "&#91;AI&#93;") {
		temp := fmt.Sprintf("[CQ:at,qq=%d] \n用户：%d\n是否开启上下文：%t\n上下文总额度：%d\n剩余额度：%d\n输入带有[设定]可以diy回复\n回复[开启上下文]启用聊天记忆\n回复[关闭上下文]停用聊天记忆\n回复[清空上下文]重置聊天", body.Sender.UserID, body.Sender.UserID, user.EnableContext, user.MaxContexts, int64(user.MaxContexts)-count)
		replyContent := Models.SendGroupMessage{
			GroupID: body.GroupID,
			Message: temp,
		}
		common.PostToCQHTTPNoResponse(replyContent, "/send_group_msg")
		return
	}

	if strings.Contains(body.Message, "&#91;清空上下文&#93;") {
		Models.DB.Model(&Models.ChatGPTContext{}).Where("user = ? and state = ?", user.User, "enable").Update("state", "disable")
		common.PostToCQHTTPNoResponse(Models.SendGroupMessage{
			GroupID: body.GroupID,
			Message: "[CQ:at,qq=" + strconv.FormatInt(body.Sender.UserID, 10) + "] 清除完了哦",
		}, "/send_group_msg")
		return
	}

	if strings.Contains(body.Message, "&#91;开启上下文&#93;") {
		user.EnableContext = true
		lock.Lock()
		Models.ChatGPTUsers[body.Sender.UserID] = user
		lock.Unlock()
		Models.DB.Save(&user)
		common.PostToCQHTTPNoResponse(Models.SendGroupMessage{
			GroupID: body.GroupID,
			Message: "[CQ:at,qq=" + strconv.FormatInt(body.Sender.UserID, 10) + "] 开启了哦",
		}, "/send_group_msg")
		return
	}

	if strings.Contains(body.Message, "&#91;关闭上下文&#93;") {
		user.EnableContext = false
		lock.Lock()
		Models.ChatGPTUsers[body.Sender.UserID] = user
		lock.Unlock()
		Models.DB.Save(&user)
		common.PostToCQHTTPNoResponse(Models.SendGroupMessage{
			GroupID: body.GroupID,
			Message: "[CQ:at,qq=" + strconv.FormatInt(body.Sender.UserID, 10) + "] 关闭了哦",
		}, "/send_group_msg")
		return
	}

	lock.RLock()
	ownerID := config.Config.General.Owner
	lock.RUnlock()

	if strings.Contains(body.Message, "&#91;调整额度&#93;") && body.Sender.UserID == ownerID {
		reg := regexp.MustCompile("&#91;[0-9]+-[0-9]+&#93;")
		if reg == nil {
			common.Logln(2, "正则匹配失败")
			return
		}
		temp := reg.FindAllStringSubmatch(body.Message, -1)
		if len(temp) == 0 {
			common.ErrorResponse(true, body.GroupID, errors.New("格式有误，应为“[调整额度] [QQID-额度]”"))
			return
		}
		data := temp[0][0]
		dataList := strings.Split(data, "-")
		UserID, err := strconv.ParseInt(dataList[0][5:], 10, 64)
		if err != nil {
			common.Logln(2, err)
			return
		}
		quota, err := strconv.Atoi(dataList[1][:len(dataList[1])-5])
		if err != nil {
			common.Logln(2, err)
			return
		}

		lock.Lock()
		target := Models.ChatGPTUsers[UserID]
		target.MaxContexts = quota
		Models.ChatGPTUsers[UserID] = target
		lock.Unlock()

		Models.DB.Save(&target)
		common.PostToCQHTTPNoResponse(Models.SendGroupMessage{
			GroupID: body.GroupID,
			Message: fmt.Sprintf("[CQ:at,qq=%d] 用户%v的额度已调整为%d", body.Sender.UserID, UserID, quota),
		}, "/send_group_msg")
		return
	}

	if user.EnableContext {
		ChatWithContext(body, user)
	} else {
		ChatGPT(body, "user")
	}
}
