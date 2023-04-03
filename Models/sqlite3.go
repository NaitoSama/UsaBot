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
	var user []ChatGPTUserInfo
	result := DB.Find(&user)
	if result.Error != nil {
		panic("failed to find user info")
	}
	for _, v := range user {
		ChatGPTUsers[v.User] = v
	}
}
