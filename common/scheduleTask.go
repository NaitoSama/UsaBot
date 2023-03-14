package common

import (
	"github.com/jasonlvhit/gocron"
)

var ScheduleClient *gocron.Scheduler

func init() {
	go Init()
}

func Init() {
	c := gocron.NewScheduler()
	ScheduleClient = c
	<-c.Start()
}

//func init() {
//	c := cron.New()
//	ScheduleClient = c
//}
//
//// ScheduleTask 添加一个定时任务
//func ScheduleTask(schedule string, function func()) (int, error) {
//	entryID, err := ScheduleClient.AddFunc(schedule, function)
//	if err != nil {
//		return -1, err
//	}
//	return int(entryID), nil
//}
