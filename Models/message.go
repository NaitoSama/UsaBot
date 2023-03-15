package Models

type Message struct {
	Font        int    `json:"font,omitempty"`
	PostType    string `json:"post_type"` // message-消息 request-请求 notice-通知 meta_event-元事件
	MessageType string `json:"message_type,omitempty"`
	GroupID     int64  `json:"group_id,omitempty"` // 私聊-临时会话发起的群号 群聊-群号
	Message     string `json:"message,omitempty"`
	RawMessage  string `json:"raw_message,omitempty"`
	SelfID      int64  `json:"self_id,omitempty"`
	Sender      sender `json:"sender,omitempty"`
	SubType     string `json:"sub_type,omitempty"`
	TargetID    int64  `json:"target_id,omitempty"`
	Time        int64  `json:"time,omitempty"`
	UserID      int64  `json:"user_id,omitempty"`
}

type sender struct {
	NickName string `json:"nickname"`
	UserID   int64  `json:"user_id"`        // 发送者ID
	Role     string `json:"role,omitempty"` // 群聊中的发送者的身份
	Card     string `json:"card,omitempty"` // 群聊中发送者的群名片
}

type SendPrivateMessage struct {
	UserID     int64  `json:"user_id"`
	GroupID    int64  `json:"group_id,omitempty"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape,omitempty"` // 是否以纯文本发送，只在message为string有效
}

type SendGroupMessage struct {
	GroupID    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape,omitempty"` // 是否以纯文本发送，只在message为string有效
}

type SendGroupMessageResponse struct {
	Data    SendGroupMessageResponseData `json:"data,omitempty"`
	RetCode int                          `json:"retcode"`
	Status  string                       `json:"status"`            // 执行成功-ok 执行失败-failed
	Msg     string                       `json:"msg,omitempty"`     // 错误信息
	Wording string                       `json:"wording,omitempty"` // 具体错误信息
}

type SendGroupMessageResponseData struct {
	MessageID int64 `json:"message_id"`
}
