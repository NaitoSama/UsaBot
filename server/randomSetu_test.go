package server

import (
	"fmt"
	"log"
	"regexp"
	"testing"
)

func Test001(t *testing.T) {
	content := "来点凯露色图"

	regl := regexp.MustCompile("来点.*[色涩瑟]图")
	if regl == nil {
		log.Println("正则解析失败")
		return
	}
	result := regl.FindAllStringSubmatch(content, -1)
	if len(result) == 0 {
		return
	}
	resultS := result[0][0]
	fmt.Println(resultS)
	fmt.Println(len(resultS))
	if len(resultS) == 12 {
		return
	}
	fmt.Println(resultS[6 : len(resultS)-6])

}
