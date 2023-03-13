package common

import "time"

// GetDayInYear 获取当前是年份第几天
func GetDayInYear(year int) int {
	monthList := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	//var leapYear bool
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		//leapYear = true
		monthList[1] = 29
	}
	month := int(time.Now().Month())
	day := time.Now().Day()
	days := 0
	for k, v := range monthList {
		if k < month-1 {
			days += v
		} else {
			days += day
			return days
		}
	}
	return -1
}

// GetAllDaysByYear 当前年份的所有天数
func GetAllDaysByYear(year int) int {
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		return 366
	} else {
		return 365
	}
}
