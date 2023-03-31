package config

import (
	"log"
	"testing"
)

func TestConfig(t *testing.T) {
	log.Println(Config.PixivPicGetter.PixivProxy)
}
