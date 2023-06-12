package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

func MoyuRili(groupIDs []int64) {
	content, err := MoyuRiliMainHandler()
	if err != nil {
		common.Logln(2, err)
		return
	}
	for _, v := range groupIDs {
		_, err = common.GroupChatSender(v, content)
		if err != nil {
			common.Logln(2, err)
		}
	}
}

func MoyuRiliMainHandler() (string, error) {
	url := "https://api.emoao.com/api/moyu?type=json"
	picApi := "https://api.vvhan.com/api/moyu?type=json"
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	respBody := Models.MoyuRili{}
	err = json.Unmarshal(respData, &respBody)
	if err != nil {
		return "", errors.New("获取摸鱼日历接口失败，无法解读响应")
	}
	if respBody.Code != 200 {
		return "", errors.New("获取摸鱼日历接口失败，响应代码不为200")
	}
	now := time.Now()
	dataTime := now.Format("2006-01-02")
	picPath := "./pic/MoyuRili-" + dataTime + ".png"

	picResp, err := common.RequestTo(picApi, "GET", "", nil)
	if err != nil {
		return "", err
	}
	defer picResp.Body.Close()
	picBody := Models.MoyuRiliPic{}
	picData, err := io.ReadAll(picResp.Body)
	if err != nil {
		return "", err
	}
	if err = json.Unmarshal(picData, &picBody); err != nil || !picBody.Success {
		return "", errors.New("摸鱼日历图片接口获取失败")
	}

	err = common.DownloadPic(picPath, picBody.Url)
	if err != nil {
		return "", err
	}
	picBase, err := common.PicBase64(picPath)
	if err != nil {
		return "", err
	}
	content := respBody.Title + "\n" + picBase
	return content, nil
}
