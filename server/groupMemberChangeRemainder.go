package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"encoding/json"
	"fmt"
)

// GroupIncrease 新群员提醒
func GroupIncrease(data []byte) {
	nts := Models.EventGMIncrease{}
	ntsData, err := json.Marshal(data)
	if err != nil {
		common.Logln(2, err)
		return
	}
	err = json.Unmarshal(ntsData, &nts)
	if err != nil {
		common.Logln(2, err)
		return
	}
	_, err = common.GroupChatSender(nts.GroupID, common.At(nts.UserID)+"欢迎新大佬，群地位-1")
	if err != nil {
		common.Logln(2, err)
	}
}

// GroupDecrease 退群提醒
func GroupDecrease(data []byte) {
	nts := Models.EventGMDecrease{}
	ntsData, err := json.Marshal(data)
	if err != nil {
		common.Logln(2, err)
		return
	}
	err = json.Unmarshal(ntsData, &nts)
	if err != nil {
		common.Logln(2, err)
		return
	}
	_, err = common.GroupChatSender(nts.GroupID, fmt.Sprintf("(%d)退群了", nts.UserID))
	if err != nil {
		common.Logln(2, err)
	}
}
