package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func RandomFood(msg Models.Message) {
	err := mainHandler(msg)
	if err != nil {
		common.ErrorResponse(true, msg.GroupID, err)
		common.Logln(2, err)
		return
	}
}

func mainHandler(msg Models.Message) error {
	files, err := common.GetFilesInfo("./pic/randomFood")
	if err != nil {
		return err
	}
	fileNum := len(files)
	randNum := common.RandIntn(fileNum)
	file := files[randNum]

	picFile, err := os.Open("./pic/randomFood/" + file.Name())
	if err != nil {
		return err
	}
	picData, err := io.ReadAll(picFile)
	if err != nil {
		return err
	}
	picBase64 := base64.StdEncoding.EncodeToString(picData)

	content := fmt.Sprintf("[CQ:image,file=base64://%s]\n[CQ:at,qq=%d] 今天吃%s", picBase64, msg.Sender.UserID, file.Name()[:len(file.Name())-4])
	replyContent := Models.SendGroupMessage{
		GroupID: msg.GroupID,
		Message: content,
	}
	common.PostToCQHTTPNoResponse(replyContent, "/send_group_msg")
	return nil
}
