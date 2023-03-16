package common

import (
	"UsaBot/Models"
	"log"
	"testing"
)

func TestReadJSON(t *testing.T) {
	boss := &Models.BossConfig{}
	err := ReadJSON("F:\\GoCode\\UsaBot\\config\\boss.json", boss)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(boss.BossHP.One.FourthRound)
}
