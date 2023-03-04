package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

const url = "http://127.0.0.1:5700"
const accessToken = ""
const contentType = "application/json"
const authorization = ""

func Tester(c *gin.Context) {
	data, _ := c.GetRawData()
	body := make(map[string]interface{})
	_ = json.Unmarshal(data, &body)
	log.Println(body)
	content := fmt.Sprintln(body)
	common.TXTWriter(content)
}

func MainHandle(c *gin.Context) {
	body := Models.Message{}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		log.Printf("Error:%s\n", err.Error())
		return
	}
	if body.PostType == "meta_event" {
		return
	}
	log.Println(body)
	content := fmt.Sprintln(body)
	common.TXTWriter(content)
	echo(body)
}

func echo(body Models.Message) {
	content := Models.SendMessage{
		UserID:  body.Sender.UserID,
		Message: body.Message,
	}
	configData, _ := json.Marshal(content)
	param := bytes.NewBuffer(configData)
	client := http.DefaultClient
	req, err := http.NewRequest("POST", url+"/send_private_msg", param)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", authorization)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	result, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(result))
}
