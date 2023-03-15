package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"encoding/base64"
	"encoding/json"
	"errors"
	_ "image/png"
	"io"
	"log"
	"os"
	"regexp"
)

// PixivPicGetter 提取pixiv图片 通过pixiv.cat 代理下载到本地，base64编码后对cqhttp发送请求
func PixivPicGetter(msg Models.Message) {
	message := Models.SendGroupMessage{
		GroupID:    msg.GroupID,
		Message:    "下载中...（如果没有图片那就是被马叔叔吃了）",
		AutoEscape: false,
	}
	common.PostToCQHTTPNoResponse(message, "/send_group_msg")
	pid, err := PixivPicID(msg.Message)
	if err != nil {
		log.Println(err)
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}
	// 图片下载
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}
	picPath := pwd + "/pic/" + pid + ".png"

	_, err = os.Stat(picPath)
	if err != nil {
		url := "https://pixiv.cat/" + pid + ".png"
		err = common.DownloadPicWithProxy(picPath, url)
		if err != nil {
			log.Println(err)
			common.ErrorResponse(true, msg.GroupID, err)
			return
		}
	}

	// 发送
	err = readPicAndSend(picPath, msg, pid)
	if err != nil {
		log.Println(err)
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}

}

// PixivPicID 获取pixiv id
func PixivPicID(message string) (string, error) {
	regl := regexp.MustCompile(`pid[0-9]+`)
	if regl == nil {
		log.Println("正则解析失败")
		return "", errors.New("正则解析失败")
	}
	result := regl.FindAllStringSubmatch(message, -1)
	if len(result) == 0 {
		return "", errors.New("请求中格式不正确，应为pidxxxxx（x为数字）")
	}
	return result[0][0][3:], nil
}

// readPicAndSend base64编码发送信息
func readPicAndSend(picPath string, msg Models.Message, pid string) error {
	picFile, err := os.Open(picPath)
	if err != nil {
		return err
	}

	picData, err := io.ReadAll(picFile)
	if err != nil {
		return err
	}

	picBase64 := base64.StdEncoding.EncodeToString(picData)

	//content := "pid:" + pid + "\n[CQ:image,file=https://pixiv.re/" + pid + ".png]"
	content := "pid:" + pid + "\n[CQ:image,file=base64://" + picBase64 + "]"
	message := Models.SendGroupMessage{
		GroupID:    msg.GroupID,
		Message:    content,
		AutoEscape: false,
	}
	response, err := common.PostToCQHTTPWithResponse(message, "/send_group_msg")
	if err != nil {
		return err
	}
	defer response.Body.Close()
	respData, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	respStruct := Models.SendGroupMessageResponse{}
	err = json.Unmarshal(respData, &respStruct)

	// 发送失败判断
	if respStruct.Status != "ok" {
		message = Models.SendGroupMessage{
			GroupID: msg.GroupID,
			Message: "涩图太涩捏，发不出来，错误信息：" + respStruct.Wording,
		}
		common.PostToCQHTTPNoResponse(message, "/send_group_msg")
	}

	return nil
}

//func downloadPicAndSend(pid string, picPath string, msg Models.Message) error {
//	url := "https://pixiv.cat/" + pid + ".png"
//	err := common.DownloadPic(picPath, url)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//
//	//content := "pid:" + pid + "\n[CQ:image,file=https://pixiv.re/" + pid + ".png]"
//	content := "pid:" + pid + "\n[CQ:image,file://" + picPath + "]"
//	log.Println(content)
//	message := Models.SendGroupMessage{
//		GroupID:    msg.GroupID,
//		Message:    content,
//		AutoEscape: false,
//	}
//	common.PostToCQHTTPNoResponse(message, "/send_group_msg")
//	return nil
//}
