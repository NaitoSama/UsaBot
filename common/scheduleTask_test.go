package common

import (
	"log"
	"testing"
	"time"
)

func TestScheduleTask(t *testing.T) {
	err := ScheduleClient.Every(1).Second().Do(func() {
		log.Print(".")
	})
	if err != nil {
		log.Println(err)
		return
	}
	for true {
		time.Sleep(time.Second)
	}
}
