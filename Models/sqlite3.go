package Models

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open("./chatGPT_context.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return
	}
	err = db.AutoMigrate(&ChatGPTContext{}, &ChatGPTUserInfo{})
	if err != nil {
		panic("failed to migrate database")
		return
	}
	DB = db
}
