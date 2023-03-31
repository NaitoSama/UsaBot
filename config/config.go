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
	DailyNews
	HolidayRemainder
	PixivPicGetter
	RandomSetu
	Soutu
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
