package common

import (
	"UsaBot/Models"
	"fmt"
)

func ErrorResponse(isGroup bool, targetID int64, err error) {
	msgContent := fmt.Sprintf("内部错误：%s", err.Error())
	if isGroup {
		content := Models.SendGroupMessage{
			GroupID: targetID,
			Message: msgContent,
		}
		PostToCQHTTPNoResponse(content, "/send_group_msg")
	} else {
		content := Models.SendPrivateMessage{
			UserID:  targetID,
			Message: msgContent,
		}
		PostToCQHTTPNoResponse(content, "/send_private_msg")
	}
}
