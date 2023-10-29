package server

import (
	"UsaBot/Models"
	"UsaBot/common"
	"UsaBot/config"
	"encoding/json"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
)

// RandomSetu Get a random anime pic from Lolicon
func RandomSetu(msg Models.Message) {
	setuConfig := config.Config
	message1 := Models.SendGroupMessage{
		GroupID: msg.GroupID,
		Message: "下载中...（如果没有发图那就是被马叔叔吃了）",
	}
	common.PostToCQHTTPNoResponse(message1, "/send_group_msg")
	tag, err := parseMsg(msg)
	if err != nil {
		common.Logln(2, err)
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}
	// Get pic id from Lolicon's Api
	setu, err := requestLoliconApi(tag)
	if err != nil {
		common.Logln(2, err)
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}
	setuData := setu.Data[0]
	// Download pic
	pwd, err := os.Getwd()
	if err != nil {
		common.Logln(2, err)
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}
	picPath := pwd + "/pic/" + strconv.FormatInt(setuData.PID, 10) + ".png"
	// Download pic with proxy or not
	if setuConfig.RandomSetu.UseProxy {
		err = common.DownloadPicWithProxy(picPath, setuData.Urls.Original)
	} else {
		err = common.DownloadPic(picPath, setuData.Urls.Original)
	}

	if err != nil {
		common.Logln(2, err)
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}

	picCQ, err := common.PicBase64(picPath)
	if err != nil {
		common.Logln(2, err)
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}

	content := "标题：" + setuData.Title + "\n作者：" + setuData.Author + "\nPID：" + strconv.FormatInt(setuData.PID, 10) + "\n原图网址：" + setuData.Urls.Original
	//content := "[CQ:image,file=" + setuData.Urls.Original + "]\n标题：" + setuData.Title + "\n作者：" + setuData.Author + "\nPID：" + strconv.FormatInt(setuData.PID, 10) + "\n原图网址：" + setuData.Urls.Original
	//content := picCQ + "\n标题：" + setuData.Title + "\n作者：" + setuData.Author + "\nPID：" + strconv.FormatInt(setuData.PID, 10) + "\n原图网址：" + setuData.Urls.Original
	message := Models.SendGroupMessage{
		GroupID: msg.GroupID,
		Message: content + "\n" + picCQ,
	}
	response, err := common.PostToCQHTTPWithResponse(message, "/send_group_msg")
	if err != nil {
		common.Logln(2, err)
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}
	defer response.Body.Close()
	respData, err := io.ReadAll(response.Body)
	if err != nil {
		common.Logln(2, err)
		common.ErrorResponse(true, msg.GroupID, err)
		return
	}

	respStruct := Models.SendGroupMessageResponse{}
	err = json.Unmarshal(respData, &respStruct)

	// can not send msg
	if respStruct.Status != "ok" {
		message = Models.SendGroupMessage{
			GroupID: msg.GroupID,
			Message: "涩图太涩捏，发不出来，错误信息：" + respStruct.Wording + "\n将以特殊方法发出",
		}
		common.PostToCQHTTPNoResponse(message, "/send_group_msg")
		err := specialSendPic(picPath, msg.GroupID, content)
		if err != nil {
			message = Models.SendGroupMessage{
				GroupID: msg.GroupID,
				Message: "涩图太涩惹，特殊形式也发不出捏",
			}
			common.PostToCQHTTPNoResponse(message, "/send_group_msg")
		}
	}
}

// parseMsg 解析msg 如果有tag返回，否则返回空字符串
func parseMsg(msg Models.Message) (string, error) {
	regl := regexp.MustCompile("来点.*[色涩瑟]图")
	if regl == nil {
		common.Logln(2, "正则解析失败")
		return "", errors.New("正则解析失败")
	}
	result := regl.FindAllStringSubmatch(msg.Message, -1)
	if len(result) == 0 {
		return "", errors.New("请求中格式不正确，应为来点<tag>色图")
	}
	resultS := result[0][0]
	if len(resultS) == 12 {
		return "", nil
	}
	return resultS[6 : len(resultS)-6], nil
}

// requestLoliconApi 请求loliconapi获取涩图信息
func requestLoliconApi(tag string) (*Models.LoliconApiResp, error) {
	reqJson := Models.LoliconApi{
		R18:       2,
		Num:       1,
		Tag:       nil,
		ExcludeAI: true,
	}

	if len(tag) != 0 {
		tagList := make([]string, 0)
		tagList = append(tagList, tag)
		reqJson.Tag = tagList
	}
	resp, err := common.RequestTo("https://api.lolicon.app/setu/v2", "POST", "application/json", reqJson)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	respStruct := &Models.LoliconApiResp{}
	err = json.Unmarshal(data, respStruct)
	if err != nil {
		return nil, err
	}
	return respStruct, nil
}

func specialSendPic(picPath string, groupNum int64, content string) error {
	newPicPath, err := common.ModifyPicMD5(picPath)
	if err != nil {
		return err
	}
	picCQCode, err := common.PicBase64(newPicPath)
	if err != nil {
		return err
	}
	message := Models.SendGroupMessage{
		GroupID: groupNum,
		Message: content + "\n" + picCQCode,
	}
	_, err = common.PostToCQHTTPWithResponse(message, "/send_group_msg")
	if err != nil {
		return err
	}
	return nil
}
