package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"encoding/base64"
	"errors"
	_ "image/png"
	"io"
	"log"
	"os"
	"regexp"
)

func PixivPicGetter(msg Models.Message) {
	pid, err := PixivPicID(msg.Message)
	if err != nil {
		log.Println(err)
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}

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
		err = common.DownloadPic(picPath, url)
		if err != nil {
			log.Println(err)
			common.ErrorResponse(true, msg.GroupID, err)
			return
		}
	}

	err = readPicAndSend(picPath, msg, pid)
	if err != nil {
		log.Println(err)
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}

}

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
	common.PostToCQHTTPNoResponse(message, "/send_group_msg")
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
