package Models

type MoyuRili struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Title   string `json:"title"`
	ImgUrl  string `json:"imgurl"`
}

type MoyuRiliPic struct {
	Success bool   `json:"success"`
	Url     string `json:"url"`
}
