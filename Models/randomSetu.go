package Models

type LoliconApi struct {
	R18       int      `json:"r18,omitempty"`       // 0-noH 1-H 2-mix default 0
	Num       int      `json:"num,omitempty"`       // 请求图片数量 default 1
	UID       []int    `json:"uid,omitempty"`       // 指定作者
	Keyword   string   `json:"keyword,omitempty"`   // 关键字
	Tag       []string `json:"tag,omitempty"`       // 通过tag过滤 如 ["pcr","黑发"]
	Size      []string `json:"size,omitempty"`      // 图片大小 default original
	Proxy     string   `json:"proxy,omitempty"`     // 代理 default i.pixiv.re
	ExcludeAI bool     `json:"excludeAI,omitempty"` // 是否排除AI default false
}

type LoliconApiResp struct {
	Error string               `json:"error"`
	Data  []LoliconApiRespData `json:"data,omitempty"`
}

type LoliconApiRespData struct {
	PID        int64                 `json:"pid"`
	P          int                   `json:"p"` // 第几张图
	UID        int64                 `json:"uid"`
	Title      string                `json:"title"`
	Author     string                `json:"author"`
	R18        bool                  `json:"r18"`
	Width      int                   `json:"width"`
	Height     int                   `json:"height"`
	Tag        []string              `json:"tag,omitempty"`
	Ext        string                `json:"ext"`        // 扩展名
	AiType     int                   `json:"aiType"`     // is ai? 0-unknown 1-no 2-yes
	UploadDate int                   `json:"uploadDate"` // 上传时间戳
	Urls       LoliconApiRespDataUrl `json:"urls"`
}

type LoliconApiRespDataUrl struct {
	Original string `json:"original,omitempty"`
	Regular  string `json:"regular,omitempty"`
	Small    string `json:"small,omitempty"`
	Thumb    string `json:"thumb,omitempty"`
	Mini     string `json:"mini,omitempty"`
}
