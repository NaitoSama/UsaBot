package Models

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var PcrDB *gorm.DB

type Boss struct {
	BossID    int
	BossValue int64
	BossStage int
	BossRound int
	GuaShu    string
	JinRu     string
}

type BossConfig struct {
	BossHP    `json:"boss_HP"`
	BossStage `json:"boss_stage"`
}

type BossHP struct {
	One   BossID `json:"one"`
	Two   BossID `json:"two"`
	Three BossID `json:"three"`
	Four  BossID `json:"four"`
	Five  BossID `json:"five"`
}

type BossID struct {
	FirstStage  int64 `json:"first_stage"`
	SecondStage int64 `json:"second_stage"`
	ThirdStage  int64 `json:"third_stage"`
	FourthStage int64 `json:"fourth_stage"`
	FifthStage  int64 `json:"fifth_stage"`
}

type BossStage struct {
	FirstStage  string `json:"first_stage"`
	SecondStage string `json:"second_stage"`
	ThirdStage  string `json:"third_stage"`
	FourthStage string `json:"fourth_stage"`
	FifthStage  string `json:"fifth_stage"`
}

func init() {
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: time.Second,
		LogLevel:      logger.Info,
		Colorful:      true,
	})

	db, err := gorm.Open(sqlite.Open("pcr_clan_battle.db"), &gorm.Config{Logger: newLogger})
	if err != nil {
		text := fmt.Sprintf("failed to connect mysql, err:%s", err.Error())
		log.Println(text)
		return
	}
	err = db.AutoMigrate(Boss{})
	if err != nil {
		log.Println(err)
		return
	}
	PcrDB = db
}
