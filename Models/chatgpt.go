package Models

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
