package Models

type Message struct {
	Font        int    `json:"font,omitempty"`
	PostType    string `json:"post_type"` // message-消息 request-请求 notice-通知 meta_event-元事件
	MessageType string `json:"message_type,omitempty"`
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
	UserID   int64  `json:"user_id"`
}

type SendMessage struct {
	UserID     int64  `json:"user_id"`
	GroupID    int64  `json:"group_id,omitempty"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape,omitempty"` // 是否以纯文本发送，只在message为string有效
}
