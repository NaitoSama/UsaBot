package server

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestReg(t *testing.T) {
	msg := "[调整额度] [2471967424-12]"
	reg := regexp.MustCompile("&#91;[0-9]+-[0-9]+&#93;")
	if reg == nil {
		log.Println("正则匹配失败")
		return
	}
	temp := reg.FindAllStringSubmatch(msg, -1)
	if len(temp) == 0 {
		log.Println("wu")
		return
	}
	data := temp[0][0]
	dataList := strings.Split(data, "-")
	UserID, err := strconv.ParseInt(dataList[0][5:], 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	quota, err := strconv.Atoi(dataList[1][:len(dataList[1])-5])
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(UserID, "\t", quota)
}
