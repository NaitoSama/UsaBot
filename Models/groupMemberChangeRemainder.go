package Models

type EventGMIncrease struct {
	SubType    string `json:"sub_type"` // 群员加入形式 approve-通过加群验证 invite-邀请入群
	GroupID    int64  `json:"group_id"`
	OperatorID int64  `json:"operator_id,omitempty"` // 操作者ID
	UserID     int64  `json:"user_id"`               // 加群人ID
	Time       int64  `json:"time"`
	SelfID     int64  `json:"self_id"`
	PostType   string `json:"post_type"`
	NoticeType string `json:"notice_type"`
}

type EventGMDecrease struct {
	SubType    string `json:"sub_type"` // 群员退出形式 leave-主动离群 kick-被踢出群
	GroupID    int64  `json:"group_id"`
	OperatorID int64  `json:"operator_id"` // 操作者ID
	UserID     int64  `json:"user_id"`     // 离群人ID
	Time       int64  `json:"time"`
	SelfID     int64  `json:"self_id"`
	PostType   string `json:"post_type"`
	NoticeType string `json:"notice_type"`
}
