package server

import (
	"UsaBot/common"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestDailyNews(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.03c3.cn/zb/api.php", nil)
	if err != nil {
		common.Logln(2, err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		common.Logln(2, err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		common.Logln(2, err)
		return
	}
	respStruct := make(map[string]string, 0)
	err = json.Unmarshal(data, &respStruct)
	msg, ok := respStruct["msg"]
	if !ok {
		common.Logln(2, "获取接口https://api.03c3.cn/zb/api.php未有响应")
		return
	}
	if msg != "Success" {
		common.Logln(2, "获取接口https://api.03c3.cn/zb/api.php失败")
		return
	}
	imageUrl, ok := respStruct["imageUrl"]
	if !ok {
		common.Logln(2, "获取接口https://api.03c3.cn/zb/api.php图片网址失败")
		return
	}
	datatime, ok := respStruct["datatime"]
	if !ok {
		common.Logln(2, "获取接口https://api.03c3.cn/zb/api.php 数据更新时间失败")
		return
	}
	picPath := "./dailyNews-" + datatime + ".png"
	err = common.DownloadPic(picPath, imageUrl)
	if err != nil {
		common.Logln(2, err)
		return
	}
	fmt.Println("ok")
}
