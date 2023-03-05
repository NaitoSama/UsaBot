package router

import (
	"UsaBot/common"
	"UsaBot/server"
	"github.com/gin-gonic/gin"
)

func StartServer() {
	InitServer()
	r := gin.New()
	router(r)
	err := r.Run(":10086")
	if err != nil {
		common.ErrorHandle(err)
		return
	}
}

func router(r *gin.Engine) {
	r.POST("/", server.MainHandler)
}

func InitServer() {
	go server.Router()
	go server.MsgHandler()
	go server.ReqHandler()
	go server.NoticeHandler()
}
