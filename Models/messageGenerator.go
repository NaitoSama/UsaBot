package Models

type SendForwardMsg struct {
	GroupID  int64        `json:"group_id"`
	Messages []ForwardMsg `json:"messages"`
}

type ForwardMsg struct {
	Type string         `json:"type"`
	Data ForwardMsgData `json:"data"`
}

type ForwardMsgData struct {
	UserName string                  `json:"name"`
	UserID   string                  `json:"uin"`
	Content  []ForwardMsgDataContent `json:"content"`
}

type ForwardMsgDataContent struct {
	Type string                    `json:"type"`
	Data ForwardMsgDataContentData `json:"data"`
}

type ForwardMsgDataContentData struct {
	Text string `json:"text"`
}
