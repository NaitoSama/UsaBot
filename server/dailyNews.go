package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

// DailyNews 每日晨报
func DailyNews(groups []int64) {
	resp, err := common.RequestTo("https://api.03c3.cn/zb/api.php", "GET", "", nil)
	if err != nil {
		log.Println(err)
		return
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	respStruct := make(map[string]string, 0)
	err = json.Unmarshal(data, &respStruct)
	msg, ok := respStruct["msg"]
	if !ok {
		log.Println("获取接口https://api.03c3.cn/zb/api.php未有响应")
		return
	}
	if msg != "Success" {
		log.Println("获取接口https://api.03c3.cn/zb/api.php失败")
		return
	}
	imageUrl, ok := respStruct["imageUrl"]
	if !ok {
		log.Println("获取接口https://api.03c3.cn/zb/api.php图片网址失败")
		return
	}
	datatime, ok := respStruct["datatime"]
	if !ok {
		log.Println("获取接口https://api.03c3.cn/zb/api.php 数据更新时间失败")
		return
	}
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}
	picPath := pwd + "/pic/dailyNews-" + datatime + ".png"
	err = common.DownloadPic(picPath, imageUrl)
	if err != nil {
		log.Println(err)
		return
	}
	picFile, err := os.Open(picPath)
	if err != nil {
		log.Println(err)
		return
	}

	picData, err := io.ReadAll(picFile)
	if err != nil {
		log.Println(err)
		return
	}

	picBase64 := base64.StdEncoding.EncodeToString(picData)

	content := "早上好，打工人，看看前几天都发生什么事儿了\n[CQ:image,file=base64://" + picBase64 + "]"

	for _, v := range groups {
		message := Models.SendGroupMessage{
			GroupID:    v,
			Message:    content,
			AutoEscape: false,
		}
		common.PostToCQHTTPNoResponse(message, "/send_group_msg")
		time.Sleep(time.Second)
	}

}
