package config

import (
	"github.com/BurntSushi/toml"
)

var Config config

func init() {
	_, err := toml.DecodeFile("./config/config.toml", &Config)
	if err != nil {
		panic(err)
		return
	}
}

type config struct {
	General          General
	ChatGPT          ChatGPT
	DailyNews        DailyNews
	HolidayRemainder HolidayRemainder
	PixivPicGetter   PixivPicGetter
	RandomSetu       RandomSetu
	Soutu            Soutu
	PicGenerator     PicGenerator
	MoyuRili         MoyuRili
	ForwardMsgGen    ForwardMsgGen
}

type General struct {
	HttpPort  int
	CQHttpUrl string
	Proxy     string
	Owner     int64
}

type ChatGPT struct {
	Enable      bool
	Model       string
	UseProxy    bool
	Url         string
	AccessToken string
}

type DailyNews struct {
	Enable    bool
	Time      string
	GroupList []int64
}

type HolidayRemainder struct {
	Enable    bool
	Time      string
	GroupList []int64
}

type PixivPicGetter struct {
	Enable     bool
	PixivProxy string
	UseProxy   bool
}

type RandomSetu struct {
	Enable   bool
	UseProxy bool
}

type Soutu struct {
	Enable         bool
	SaucenaoApi    string
	SaucenaoApiKey string
	Results        int
}

type PicGenerator struct {
	Enable      bool
	Number      int
	UseProxy    bool
	Url         string
	AccessToken string
	Size        string
}

type MoyuRili struct {
	Enable    bool
	Time      string
	GroupList []int64
}

type ForwardMsgGen struct {
	Enable bool
}
