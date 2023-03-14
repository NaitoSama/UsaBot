package router

import (
	"UsaBot/common"
	"UsaBot/server"
	"github.com/gin-gonic/gin"
)

func StartServer() {
	InitServer()
	scheduleTask()
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

func scheduleTask() {
	groupNum := []int64{1036326321, 292249427}
	common.ScheduleClient.Every(1).Day().At("14:00").Do(func() { server.HolidayReminderTask(groupNum) })
}
