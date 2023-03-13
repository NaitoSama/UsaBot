package common

import (
	"fmt"
	"testing"
	"time"
)

func TestGetWeekDay(t *testing.T) {
	weekDay := time.Now().Weekday()
	fmt.Printf("week day is %d, type is %T", weekDay, weekDay)
}

func TestGetDayInYear(t *testing.T) {
	fmt.Println(GetDayInYear(2023))
}
