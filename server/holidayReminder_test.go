package server

import (
	"fmt"
	"log"
	"testing"
)

func TestHolidayReminder(t *testing.T) {
	result, err := HolidayReminder()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(result)
}
