package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"UsaBot/config"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func BoardCast(msg Models.Message) {
	lock.RLock()
	masterID := config.Config.General.Owner
	lock.RUnlock()
	if msg.Sender.UserID != masterID {
		common.PrivateChatSender(msg.Sender.UserID, "您暂时没有权限使用广播功能哦")
		return
	}
	failedList, failedNum, err := boardCastLogic(msg)
	if err != nil {
		common.ErrorResponse(false, msg.Sender.UserID, err)
		return
	}
	err = common.PrivateChatSender(msg.Sender.UserID, "发送完成，其中"+strconv.Itoa(failedNum)+"条发送失败，分别是:"+fmt.Sprint(failedList))
	if err != nil {
		common.ErrorResponse(false, msg.Sender.UserID, err)
		return
	}
}

func boardCastLogic(msg Models.Message) ([][2]string, int, error) {
	var (
		groupIDs = make([]int64, 0)
		//token string
		content    string
		failedNum  int
		failedList [][2]string
	)
	msgList := strings.Split(msg.Message, " ")
	for _, v := range msgList {
		if strings.Contains(v, "groupID") {
			vList := strings.Split(v, "=")
			groupIDList := strings.Split(vList[len(vList)-1], ",")
			for _, s := range groupIDList {
				groupID, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return nil, failedNum, errors.New("groupID格式不正确，应为“groupID=12345,54321,1234567”")
				}
				groupIDs = append(groupIDs, groupID)
			}
		} else if strings.Contains(v, "content") {
			vList := strings.Split(v, "=")
			content = vList[len(vList)-1]
			if content == "" || strings.Contains(content, "content") {
				return nil, failedNum, errors.New("content格式错误，应为“content=在吗？”")
			}
		}
	}

	for _, v := range groupIDs {
		_, err := common.GroupChatSender(v, content)
		if err != nil {
			failedNum++
			temp := [2]string{strconv.FormatInt(v, 10), err.Error()}
			failedList = append(failedList, temp)
		}
	}
	return failedList, failedNum, nil
}
