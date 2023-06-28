package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
)

func MsgGenMain(body Models.Message) {
	result, err := msgExtract(body.Message, body.GroupID)
	if err != nil {
		common.ErrorResponse(true, body.GroupID, err)
		return
	}
	common.PostToCQHTTPNoResponse(result, "/send_group_forward_msg")
}

func msgExtract(msg string, groupID int64) ([]Models.ForwardMsg, error) {
	msgList := strings.Split(msg, "|")
	if len(msgList) <= 1 {
		return nil, errors.New("格式不正确，应为“[QQ号]:正文|[QQ号]:正文”")
	}
	result := make([]Models.ForwardMsg, len(msgList))
	reql := regexp.MustCompile("[0-9]{5,}.*:.*")
	if reql == nil {
		return nil, errors.New("生成消息正则解析失败")
	}
	for i := 0; i < len(msgList); i++ {
		match := reql.FindAllStringSubmatch(msgList[i], -1)
		if len(match) == 0 {
			return nil, errors.New("格式不正确，应为“[QQ号]:正文|[QQ号]:正文”")
		}
		matchContent := match[0][0]
		index := strings.Index(matchContent, ":")
		qID, err := strconv.ParseInt(matchContent[:index], 10, 64)
		if err != nil {
			return nil, errors.New("qq号解析错误:" + err.Error())
		}
		chatMsg := matchContent[index+1:]

		payload := Models.GetQQMemberInfo{
			GroupID: groupID,
			UserID:  qID,
			NoCache: false,
		}
		resp, err := common.PostToCQHTTPWithResponse(payload, "/get_group_member_info")
		if err != nil {
			return nil, errors.New("群内没有QQ号为" + matchContent[:index] + "的成员" + err.Error())
		}
		defer resp.Body.Close()
		qqMember := Models.QQMember{}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &qqMember)
		if err != nil {
			return nil, err
		}

		forwardMsg := Models.ForwardMsg{
			Type: "node",
			Data: Models.ForwardMsgData{
				UserName: qqMember.Nickname,
				UserID:   qID,
				Content:  chatMsg,
			},
		}
		result[i] = forwardMsg
	}
}
