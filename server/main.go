package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"UsaBot/config"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"sync"
)

var (
	dataQueue = make(chan Models.Message, 10)
	msgQueue  = make(chan Models.Message, 3)
	reqQueue  = make(chan Models.Message, 3)
	ntsQueue  = make(chan Models.Message, 3)
	lock      sync.RWMutex
)

func Router() {
	for true {
		select {
		case body := <-dataQueue:
			switch body.PostType {
			case "message":
				msgQueue <- body
			case "request":
				reqQueue <- body
			case "notice":
				ntsQueue <- body
			}
		}
	}
}

func MsgHandler() {
	for true {
		select {
		case body := <-msgQueue:
			switch body.MessageType {
			case "private":
				echo(body)
			case "group":
				if strings.Contains(body.Message, "[CQ:at,qq="+strconv.FormatInt(body.SelfID, 10)+"]") {
					if config.Config.Soutu.Enable && strings.Contains(body.Message, "搜图") {
						souTu(body)
					} else if config.Config.PixivPicGetter.Enable && strings.Contains(body.Message, "提取图片") {
						PixivPicGetter(body)
					} else if config.Config.RandomSetu.Enable && strings.Contains(body.Message, "来点") {
						RandomSetu(body)
					} else if config.Config.ChatGPT.Enable {

						lock.RLock()
						user, ok := Models.ChatGPTUsers[body.Sender.UserID]
						if !ok {
							user = Models.ChatGPTUserInfo{
								User:          body.Sender.UserID,
								EnableContext: false,
								MaxContexts:   50,
							}
							Models.ChatGPTUsers[body.Sender.UserID] = user
							Models.DB.Create(&user)
						}
						lock.RUnlock()

						var count int64
						Models.DB.Model(Models.ChatGPTContext{}).Where("user = ? and state = ?", user.User, "enable").Count(&count)

						if count >= int64(user.MaxContexts) {
							common.ErrorResponse(true, body.GroupID, errors.New("您的上下文数量已达上限,现在为您进行清除"))
							Models.DB.Model(&Models.ChatGPTContext{}).Where("user = ? and state = ?", user.User, "enable").Update("state", "disable")
							count = 0
						}

						if strings.Contains(body.Message, "[AI]") {
							temp := fmt.Sprintf("[CQ:at,qq=%d] \n用户：%d\n是否开启上下文：%t\n上下文总额度：%d\n剩余额度：%d\n回复“[清空上下文]”可以重置聊天哦", body.Sender.UserID, body.Sender.UserID, user.EnableContext, user.MaxContexts, count)
							replyContent := Models.SendGroupMessage{
								GroupID: body.GroupID,
								Message: temp,
							}
							common.PostToCQHTTPNoResponse(replyContent, "/send_group_msg")
							return
						}

						if strings.Contains(body.Message, "[清空上下文]") {
							Models.DB.Model(&Models.ChatGPTContext{}).Where("user = ? and state = ?", user.User, "enable").Update("state", "disable")
							common.PostToCQHTTPNoResponse(Models.SendGroupMessage{
								GroupID: body.GroupID,
								Message: "清除完了哦",
							}, "/send_group_msg")
							return
						}

						if user.EnableContext {
							ChatWithContext(body, user)
						} else {
							ChatGPT(body, "user")
						}

						//if body.Sender.UserID == 2471967424 && strings.Contains(body.Message, "system") {
						//	temp := strings.Split(body.Message, "system")
						//	body.Message = strings.Join(temp, "")
						//	ChatGPT(body, "system")
						//} else {
						//	ChatGPT(body, "user")
						//}

					}

				}
			}
		}
	}
}

func ReqHandler() {
	for true {
		select {
		case _ = <-reqQueue:

		}
	}
}

func NoticeHandler() {
	for true {
		select {
		case _ = <-ntsQueue:

		}
	}
}

func MainHandler(c *gin.Context) {
	body := Models.Message{}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		common.Logln(2, err)
		return
	}
	dataQueue <- body

	//TestPrinter(c)
}

func TestPrinter(c *gin.Context) {
	data, _ := c.GetRawData()
	testLog := make(map[string]interface{})
	_ = json.Unmarshal(data, &testLog)
	content := fmt.Sprintln(testLog)
	common.TXTWriter(content)
}
