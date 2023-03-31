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

// todo 搜图和chatgpt还没改
type config struct {
	DailyNews        DailyNews
	HolidayRemainder HolidayRemainder
	PixivPicGetter   PixivPicGetter
	RandomSetu       RandomSetu
	Soutu            Soutu
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
}

type RandomSetu struct {
	Enable bool
}

type Soutu struct {
	Enable         bool
	SaucenaoApi    string
	SaucenaoApiKey string
	Results        int
}
