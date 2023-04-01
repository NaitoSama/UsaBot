package common

import (
	"UsaBot/Models"
	"UsaBot/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	saucenaoApi    = config.Config.Soutu.SaucenaoApi
	saucenaoApiKey = config.Config.Soutu.SaucenaoApiKey
)

func SauceNao(picUrl string, numres int) (*Models.SauceNao, error) {
	lock.RLock()
	reqUrl := fmt.Sprintf("%s?db=999&output_type=2&numres=%d&url=%s&api_key=%s", saucenaoApi, numres, picUrl, saucenaoApiKey)
	lock.RUnlock()
	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Get(reqUrl)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	result := Models.SauceNao{Header: Models.SauceNaoHeader{Status: -1}}
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &result, nil
}
