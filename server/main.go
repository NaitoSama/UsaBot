package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"UsaBot/config"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"strconv"
	"strings"
	"sync"
)

var (
	dataQueue = make(chan []byte, 10)
	msgQueue  = make(chan Models.Message, 3)
	reqQueue  = make(chan []byte, 3)
	ntsQueue  = make(chan []byte, 3)
	lock      sync.RWMutex
)

func Router() {
	for true {
		select {
		case bodyData := <-dataQueue:
			bodyMap := make(map[string]string)
			err := json.Unmarshal(bodyData, &bodyMap)
			if err != nil {
				common.Logln(2, err)
				break
			}
			switch bodyMap["post_type"] {
			case "message":
				body := Models.Message{}
				err = json.Unmarshal(bodyData, &body)
				msgQueue <- body
			case "request":
				reqQueue <- bodyData
			case "notice":
				ntsQueue <- bodyData
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
		case data := <-ntsQueue:
			ntsMap := make(map[string]string)
			err := json.Unmarshal(data, &ntsMap)
			if err != nil {
				common.Logln(2, err)
				break
			}
			ntsType, ok := ntsMap["notice_type"]
			if !ok {
				break
			}
			switch ntsType {
			case "group_increase":
				GroupIncrease(data)
			case "group_decrease":
				GroupDecrease(data)
			}
		}
	}
}

func MainHandler(c *gin.Context) {
	bodyData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		common.Logln(2, err)
	}
	//err = json.Unmarshal(bodyData, &bodyMap)
	//if err != nil {
	//	common.Logln(2, err)
	//}
	//
	//body := Models.Message{}
	//err = c.ShouldBindJSON(&body)
	//if err != nil {
	//	common.Logln(2, err)
	//	return
	//}
	dataQueue <- bodyData

	//TestPrinter(c)
}

func TestPrinter(c *gin.Context) {
	data, _ := c.GetRawData()
	testLog := make(map[string]interface{})
	_ = json.Unmarshal(data, &testLog)
	content := fmt.Sprintln(testLog)
	common.TXTWriter(content)
}
