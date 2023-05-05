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
	"strconv"
	"strings"
)

func PicGenerator(msg Models.Message) {
	s1 := strings.ReplaceAll(msg.Message, "[CQ:at,qq="+strconv.FormatInt(msg.SelfID, 10)+"]", "")
	s2 := strings.ReplaceAll(s1, "生成图片", "")
	content, err := picGeneratorMainHandler(msg.Sender.UserID, s2)
	if err != nil {
		common.ErrorResponse(true, msg.GroupID, err)
	}
	_, err = common.GroupChatSender(msg.GroupID, content)
	if err != nil {
		common.ErrorResponse(true, msg.GroupID, err)
	}
}

func picGeneratorMainHandler(userID int64, prompt string) (string, error) {
	lock.RLock()
	configData := config.Config
	lock.RUnlock()
	var client *http.Client

	body := Models.PicGeneratorReq{
		Prompt: prompt,
		Number: configData.PicGenerator.Number,
		Size:   configData.PicGenerator.Size,
	}
	data, _ := json.Marshal(body)
	param := bytes.NewBuffer(data)

	if configData.ChatGPT.UseProxy {
		proxyUrl, _ := url.Parse(configData.General.Proxy)
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		client = &http.Client{
			Transport: transport,
		}
	} else {
		client = &http.Client{}
	}

	req, err := http.NewRequest("POST", configData.PicGenerator.Url, param)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+configData.PicGenerator.AccessToken)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	resp := Models.PicGeneratorResp{}
	err = json.Unmarshal(resData, &resp)
	if err != nil {
		return "", err
	}

	content := common.At(userID)
	for k, v := range resp.Data {
		content += fmt.Sprintf("\n第%d张生成图:\n %s", k+1, common.Pic(v.Url))
	}
	return content, nil

}
