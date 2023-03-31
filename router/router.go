package router

import (
	"UsaBot/common"
	"UsaBot/config"
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

// scheduleTask 添加计划任务
func scheduleTask() {
	if config.Config.HolidayRemainder.Enable {
		common.ScheduleClient.Every(1).Day().At(config.Config.HolidayRemainder.Time).Do(func() { server.HolidayReminderTask(config.Config.HolidayRemainder.GroupList) })
	}
	if config.Config.DailyNews.Enable {
		common.ScheduleClient.Every(1).Day().At(config.Config.DailyNews.Time).Do(func() { server.DailyNews(config.Config.DailyNews.GroupList) })
	}
}
