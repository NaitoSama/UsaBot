package server

import (
	"UsaBot/Models"
	"UsaBot/common"
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
	message := msg.Message[strings.Index(msg.Message, "]")+1:]
	msgContent := Models.ChatGPTMessage{
		Role:    role,
		Content: message,
	}
	content := Models.ChatGPT{
		Model:    "gpt-3.5-turbo",
		Messages: []Models.ChatGPTMessage{msgContent},
	}
	data, _ := json.Marshal(content)
	param := bytes.NewBuffer(data)
	proxyUrl, _ := url.Parse("http://localhost:1080")
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	client := &http.Client{
		Transport: transport,
	}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", param)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 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")
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
