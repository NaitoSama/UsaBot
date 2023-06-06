package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DailyNews 每日晨报
func DailyNews(groups []int64) {
	err := DailyNewsFanXing(groups)
	if err != nil {
		common.Logln(2, err)
		return
	}
}

func DailyNews00(groups []int64) {
	//resp, err := common.RequestTo("https://api.03c3.cn/zb/api.php", "GET", "", nil)
	//if err != nil {
	//	common.Logln(2, err)
	//	return
	//}

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
	//respStruct := make(map[string]string, 0)
	respStruct := Models.DailyNews{}
	err = json.Unmarshal(data, &respStruct)
	//msg, ok := respStruct["msg"]
	//if !ok {
	//	common.Logln(2, "获取接口https://api.03c3.cn/zb/api.php未有响应")
	//	return
	//}
	//if msg != "Success" {
	//	common.Logln(2, "获取接口https://api.03c3.cn/zb/api.php失败")
	//	return
	//}
	//imageUrl, ok := respStruct["imageUrl"]
	//if !ok {
	//	common.Logln(2, "获取接口https://api.03c3.cn/zb/api.php图片网址失败")
	//	return
	//}
	//datatime, ok := respStruct["datatime"]
	//if !ok {
	//	common.Logln(2, "获取接口https://api.03c3.cn/zb/api.php 数据更新时间失败")
	//	return
	//}
	if respStruct.Code != 200 {
		common.Logln(2, "获取接口https://api.03c3.cn/zb/api.php失败")
		return
	}
	imageUrl := respStruct.ImageUrl

	now := time.Now()
	dataTime := now.Format("2006-01-02")
	picPath := "./pic/dailyNews-" + dataTime + ".png"
	err = common.DownloadPic(picPath, imageUrl)
	if err != nil {
		common.Logln(2, err)
		return
	}

	picData, err := common.PicBase64(picPath)
	if err != nil {
		common.Logln(2, err)
		return
	}
	content := fmt.Sprintf("早上好打工人，今天是%d月%d号，看看最近都发生了什么\n%s", now.Month(), now.Day(), picData)

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

// DailyNews01 每日新闻，使用“https://api.vvhan.com/60s.html”的api
func DailyNews01(groupIDs []int64) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.vvhan.com/api/60s?type=json", nil)
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

	dailyNews := Models.DailyNews_01{}
	err = json.Unmarshal(data, &dailyNews)
	if err != nil {
		temp := make(map[string]interface{})
		_ = json.Unmarshal(data, &temp)
		respBody := fmt.Sprint(temp)
		common.Logln(1, respBody)
		common.Logln(2, err)
		return
	}

	for _, v := range groupIDs {
		if !dailyNews.Success {
			common.ErrorResponse(true, v, errors.New("每日新闻请求api失败"))
		} else {
			content := fmt.Sprintf("早上好打工人，今天是%s，农历%s，%s\n为您准备了以下新闻：\n%s", dailyNews.Time[0], dailyNews.Time[1], dailyNews.Time[2], common.Pic(dailyNews.ImgUrl))
			_, err = common.GroupChatSender(v, content)
			if err != nil {
				common.ErrorResponse(true, v, errors.New("每日新闻发送失败"))
			}
		}

	}
}

func DailyNewsFanXing(groupIDs []int64) error {
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.emoao.com/api/60s?type=json", nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody := Models.DailyNewsFanXing{}
	respBodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBodyData, &respBody)
	if err != nil {
		return err
	}
	if respBody.Code != 200 {
		return errors.New("获取接口数据失败")
	}
	now := time.Now()
	dataTime := now.Format("2006-01-02")
	picPath := "./pic/dailyNews-" + dataTime + ".png"
	err = common.DownloadPic(picPath, respBody.Data.Image)
	if err != nil {
		return err
	}
	weekday := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	picBase, err := common.PicBase64(picPath)
	if err != nil {
		return err
	}
	for _, v := range groupIDs {
		content := fmt.Sprintf("早上好打工人，今天是%d月%d日，%s\n为您准备了以下新闻：\n%s", now.Month(), now.Day(), weekday[now.Weekday()], picBase)
		_, err = common.GroupChatSender(v, content)
		if err != nil {
			common.ErrorResponse(true, v, errors.New("每日新闻发送失败"))
			return err
		}
	}
	return nil
}
