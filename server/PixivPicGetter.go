package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"errors"
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
	picPath := pwd + "\\pic\\" + pid + ".png"
	log.Println(picPath)
	_, err = os.Stat(picPath)
	if err != nil {
		err1 := downloadPicAndSend(pid, picPath, msg)
		if err1 != nil {
			log.Println(err)
			common.ErrorResponse(true, msg.GroupID, err)
			return
		}
	} else {
		content := "pid:" + pid + "\n[CQ:image,file=https://pixiv.re/" + pid + ".png]"
		log.Println(content)
		message := Models.SendGroupMessage{
			GroupID:    msg.GroupID,
			Message:    content,
			AutoEscape: false,
		}
		common.PostToCQHTTPNoResponse(message, "/send_group_msg")
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

func downloadPicAndSend(pid string, picPath string, msg Models.Message) error {
	//url := "https://pixiv.cat/" + pid + ".png"
	//err := common.DownloadPic(picPath, url)
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}

	content := "pid:" + pid + "\n[CQ:image,file=https://pixiv.re/" + pid + ".png]"
	message := Models.SendGroupMessage{
		GroupID:    msg.GroupID,
		Message:    content,
		AutoEscape: false,
	}
	common.PostToCQHTTPNoResponse(message, "/send_group_msg")
	return nil
}
