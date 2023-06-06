package Models

type DailyNews_01 struct {
	Success bool     `json:"success"`
	Name    string   `json:"name"`
	Time    []string `json:"time"`
	Data    []string `json:"data"`
	ImgUrl  string   `json:"imgUrl"`
}

type DailyNews struct {
	Code     int    `json:"code"`
	Message  string `json:"msg"`
	ImageUrl string `json:"imageUrl"`
}

type DailyNewsFanXing struct {
	Code    int                  `json:"code"`
	Message string               `json:"msg,omitempty"`
	Data    DailyNewsFanXingData `json:"data,omitempty"`
	Time    int64                `json:"time,omitempty"`
	Usage   int                  `json:"usage,omitempty"`
	LogID   string               `json:"log_id,omitempty"`
}

type DailyNewsFanXingData struct {
	Date      string   `json:"date"`
	News      []string `json:"news"`
	WeiYu     string   `json:"weiyu"`
	Image     string   `json:"image"`
	HeadImage string   `json:"head_image,omitempty"`
}
