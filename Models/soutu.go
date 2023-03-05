package Models

type SauceNao struct {
	Header  SauceNaoHeader   `json:"header"`
	Results []sauceNaoResult `json:"results,omitempty"`
}

type SauceNaoHeader struct {
	Status int `json:"status"` // -1-失败 0-成功
}

type sauceNaoResult struct {
	Header sauceNaoResultHeader `json:"header"`
	Data   sauceNaoResultData   `json:"data"`
}

type sauceNaoResultHeader struct {
	Similarity string `json:"similarity"`
	Thumbnail  string `json:"thumbnail"` // 缩略图
}

type sauceNaoResultData struct {
	ExtUrls    []string `json:"ext_urls,omitempty"`    // 原图链接
	Title      string   `json:"title,omitempty"`       // 原图标题
	AuthorName string   `json:"author_name,omitempty"` // 作者名字
	AuthorUrl  string   `json:"author_url,omitempty"`  // 作者主页链接
}
