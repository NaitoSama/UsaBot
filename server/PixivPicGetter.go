package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"errors"
	"log"
	"regexp"
)

func PixivPicGetter(msg Models.Message) {
	_, err := PixivPicID(msg.Message)
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
