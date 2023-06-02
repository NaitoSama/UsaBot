package Models

type DailyNews_01 struct {
	Success bool     `json:"success"`
	Name    string   `json:"name"`
	Time    []string `json:"time"`
	Data    []string `json:"data"`
	ImgUrl  string   `json:"imgUrl"`
}
