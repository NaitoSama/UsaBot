package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// HolidayReminder 节假日定时提醒 返回消息string
func HolidayReminder() (string, error) {
	year := time.Now().Year()
	url := fmt.Sprintf("http://api.apihubs.cn/holiday/get?field=date,holiday,holiday_recess,yearday&year=%d&holiday=99&holiday_recess=1&lunar=1&order_by=1&cn=1&size=100", year)
	resp, err := common.RequestTo(url, "GET", "", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	body := Models.ApiHubs{}

	err = json.Unmarshal(bodyData, &body)
	if err != nil {
		return "", err
	}
	if body.Code != 0 {

		return "", errors.New("节假日请求失败")
	}

	holidayCN := make([]string, 0)
	holidays := make([]Models.Holidays, 0)
	for _, v := range body.Data.List {
		if len(holidayCN) != 0 && v.HolidayCN == holidayCN[len(holidayCN)-1] {
			continue
		}
		holidayCN = append(holidayCN, v.HolidayCN)
		holidays = append(holidays, Models.Holidays{
			YearDay:   v.YearDay,
			HolidayCN: v.HolidayCN,
		})
	}

	weekList := []string{"天", "一", "二", "三", "四", "五", "六"}
	now := time.Now()
	month := int(now.Month())
	day := now.Day()
	weekday := weekList[now.Weekday()]
	date := fmt.Sprintf("%d月%d日 星期%s", month, day, weekday)
	days := common.GetDayInYear(now.Year())

	contentList := make([]string, 0)
	temp := "摸鱼提醒你，今天是" + date + "。有事没事起身去茶水间，去厕所，去廊道走走。别老在工位上坐着，钱是老板的，但命是自己的。\n" + "距离周末还有" + strconv.Itoa(int(5-now.Weekday())) + "天"
	contentList = append(contentList, temp)
	for _, v := range holidays {
		if (v.YearDay - days) < 0 {
			continue
		}
		temp = "距离" + v.HolidayCN + "还有" + strconv.Itoa(v.YearDay-days) + "天"
		contentList = append(contentList, temp)
	}
	temp = "距离元旦还有" + strconv.Itoa(common.GetAllDaysByYear(now.Year())+1-days) + "\n认认真真上班，这根本就不叫赚钱，那是用劳动换取报酬。只有偷懒，在上班的时候摸鱼划水，你才是从老板手里赚到了钱。最后，祝愿天下所有摸鱼人，都能愉快的度过每一天!"
	contentList = append(contentList, temp)
	content := strings.Join(contentList, "\n")
	return content, nil
}

// HolidayReminderTask 给多个群发送节假日信息
func HolidayReminderTask(groupNum []int64) {
	msg, err := HolidayReminder()
	if err != nil {
		common.Logln(2, err)
		return
	}
	for _, v := range groupNum {
		content := Models.SendGroupMessage{
			GroupID:    v,
			Message:    msg,
			AutoEscape: false,
		}
		common.PostToCQHTTPNoResponse(content, "/send_group_msg")
		time.Sleep(time.Second)
	}
}
