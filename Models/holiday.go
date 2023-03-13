package Models

type ApiHubs struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data ApiHubsData `json:"data"`
}

type ApiHubsData struct {
	List []ApiHubsList `json:"list"`
}

type ApiHubsList struct {
	Date          int64  `json:"date"`           // 日期 如 20230122
	YearDay       int    `json:"yearday"`        // 年份天数
	HolidayRecess int    `json:"holiday_recess"` // 1-是假日 2-不是假日
	HolidayCN     string `json:"holiday_cn"`     // 节假日名称
}

type Holidays struct {
	YearDay   int    `json:"yearday"`    // 年份天数
	HolidayCN string `json:"holiday_cn"` // 节假日名称
}
