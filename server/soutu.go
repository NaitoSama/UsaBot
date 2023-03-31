package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"errors"
	"fmt"
	"strings"
)

// souTu 搜图 使用saucenao引擎
func souTu(msg Models.Message) {
	if msg.MessageType == "private" {
		common.ErrorResponse(false, msg.Sender.UserID, errors.New("暂不支持私聊搜图/汪汪"))
		return
	}
	if !strings.Contains(msg.Message, "CQ:image") {
		message := "图呢图呢？"
		common.ErrorResponse(true, msg.GroupID, errors.New(message))
		return
	}
	picUrl := msg.Message[strings.Index(msg.Message, "url=http")+4 : strings.Index(msg.Message, ";is_origin")]

	common.PostToCQHTTPNoResponse(Models.SendGroupMessage{GroupID: msg.GroupID, Message: "查询中"}, "/send_group_msg")

	result, err := common.SauceNao(picUrl, 3)
	if err != nil {
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}
	if result.Header.Status != 0 {
		common.ErrorResponse(true, msg.GroupID, errors.New("查询出错哦，是saucenao的错啦"))
		return
	}
	msgContentList := make([]string, 0)
	for k, v := range result.Results {
		temp := fmt.Sprintf("第%d个结果、\n[CQ:image,file=%s]\n相似度：%v\n原图链接：%s", k+1, v.Header.Thumbnail, v.Header.Similarity, v.Data.ExtUrls[0])
		msgContentList = append(msgContentList, temp)
		//log.Println(temp)
		//content := Models.SendGroupMessage{
		//	GroupID: msg.GroupID,
		//	Message: temp,
		//}
		//common.PostToCQHTTPNoResponse(content, "/send_group_msg")
		//time.Sleep(time.Second)
	}
	msgContent := strings.Join(msgContentList, "\n")
	content := Models.SendGroupMessage{
		GroupID: msg.GroupID,
		Message: msgContent,
	}
	common.PostToCQHTTPNoResponse(content, "/send_group_msg")
}
