package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"fmt"
	"strings"
)

func echo(body Models.Message) {
	content := Models.SendPrivateMessage{
		UserID:  body.Sender.UserID,
		Message: body.Message,
	}
	common.PostToCQHTTPNoResponse(content, "/send_private_msg")
}

func echoToGroup(msg Models.Message) {
	//if msg.Sender.Card != "" {
	//	msgSender := msg.Sender.Card
	//} else {
	//	msgSender := msg.Sender.NickName
	//}
	msgContent := strings.Join(strings.Split(msg.Message, "[CQ:at,qq=1975205178]")[:], "")
	message := fmt.Sprintf("[CQ:at,qq=%d]%s", msg.Sender.UserID, msgContent)
	content := Models.SendGroupMessage{
		GroupID: msg.GroupID,
		Message: message,
	}
	common.PostToCQHTTPNoResponse(content, "/send_group_msg")
}

//func TestHandler(c *gin.Context) {
//	body := Models.Message{}
//	err := c.ShouldBindJSON(&body)
//	if err != nil {
//		log.Printf("Error:%s\n", err.Error())
//		return
//	}
//	if body.PostType == "meta_event" {
//		return
//	}
//	log.Println(body)
//	content := fmt.Sprintln(body)
//	common.TXTWriter(content)
//	echo(body)
//}

//func Tester(c *gin.Context) {
//	data, _ := c.GetRawData()
//	body := make(map[string]interface{})
//	_ = json.Unmarshal(data, &body)
//	log.Println(body)
//	content := fmt.Sprintln(body)
//	common.TXTWriter(content)
//}
