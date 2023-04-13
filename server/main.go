package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"UsaBot/config"
	"encoding/json"
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
	lock.RLock()
	configData := config.Config
	lock.RUnlock()
	for true {
		select {
		case body := <-msgQueue:
			switch body.MessageType {
			case "private":
				echo(body)
			case "group":
				if strings.Contains(body.Message, "[CQ:at,qq="+strconv.FormatInt(body.SelfID, 10)+"]") {
					if configData.Soutu.Enable && strings.Contains(body.Message, "搜图") {
						souTu(body)
					} else if configData.PixivPicGetter.Enable && strings.Contains(body.Message, "提取图片") {
						PixivPicGetter(body)
					} else if configData.RandomSetu.Enable && strings.Contains(body.Message, "来点") {
						RandomSetu(body)
					} else if strings.Contains(body.Message, "今天吃什么") {
						RandomFood(body)
					} else if configData.ChatGPT.Enable {

						chatGPTMainHandler(body)

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
