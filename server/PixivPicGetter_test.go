package server

import (
	"UsaBot/common"
	"log"
	"testing"
)

func TestPixivPicID(t *testing.T) {
	message := "pid105998030"
	pid, err := PixivPicID(message)
	if err != nil {
		log.Println(err)
		return
	}
	picPath := "F:\\GoCode\\UsaBot\\pic\\" + pid + ".jpg"
	url := "https://pixiv.cat/" + pid + ".jpg"
	err = common.DownloadPic(picPath, url)
	if err != nil {
		log.Println(err)
		return
	}
}
