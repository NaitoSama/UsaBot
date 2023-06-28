package Models

type ForwardMsg struct {
	Type string         `json:"type"`
	Data ForwardMsgData `json:"data"`
}

type ForwardMsgData struct {
	UserName string `json:"name"`
	UserID   int64  `json:"uin"`
	Content  string `json:"content"`
}
