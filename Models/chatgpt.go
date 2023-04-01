package Models

import "gorm.io/gorm"

var ChatGPTUsers = make(map[int64]ChatGPTUserInfo)

func init() {
	var user []ChatGPTUserInfo
	result := DB.Find(&user)
	if result.Error != nil {
		panic("failed to find user info")
	}
	for _, v := range user {
		ChatGPTUsers[v.User] = v
	}
}

type ChatGPT struct {
	Model    string           `json:"model"` // 模型 一般为 “gpt-3.5-turbo”
	Messages []ChatGPTMessage `json:"messages"`
}

type ChatGPTMessage struct {
	Role    string `json:"role"` // 角色 有 user assistant system
	Content string `json:"content"`
}

type ChatGPTResponse struct {
	ID      string          `json:"id"`
	Object  string          `json:"object"`
	Created int64           `json:"created"`
	Model   string          `json:"model"`
	Usage   ChatGPTUsage    `json:"usage"`
	Choices []ChatGPTChoice `json:"choices"`
}

type ChatGPTUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatGPTChoice struct {
	Message      ChatGPTMessage `json:"message"`
	FinishReason string         `json:"finish_reason"`
	Index        int            `json:"index"`
}

type ChatGPTContext struct {
	gorm.Model
	Role    string
	Content string
	User    int64
	State   string
}

func (*ChatGPTContext) TableName() string {
	return "chatgpt_context"
}

type ChatGPTUserInfo struct {
	gorm.Model
	User          int64
	EnableContext bool
	MaxContexts   int
}

func (*ChatGPTUserInfo) TableName() string {
	return "chatgpt_user_info"
}
