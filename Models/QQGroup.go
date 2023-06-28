package Models

type QQMember struct {
	GroupID  int64  `json:"group_id"`
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
}

type GetQQMemberInfo struct {
	GroupID int64 `json:"group_id"`
	UserID  int64 `json:"user_id"`
	NoCache bool  `json:"no_cache"`
}
