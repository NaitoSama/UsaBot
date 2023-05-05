package Models

type PicGeneratorReq struct {
	Prompt string `json:"prompt"`
	Number int    `json:"n"`
	Size   string `json:"size"` // 尺寸 如: 1024x1024 1920x1080
}

type PicGeneratorResp struct {
	Created int64              `json:"created"`
	Data    []PicGeneratorData `json:"data"`
}

type PicGeneratorData struct {
	Url string `json:"url"`
}
