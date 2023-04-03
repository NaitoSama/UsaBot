package Models

import (
	"UsaBot/common"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open("./chatGPT_context.db"), &gorm.Config{})
	if err != nil {
		common.Logln(2, "failed to connect database")
		panic("failed to connect database")
		return
	}
	err = db.AutoMigrate(&ChatGPTContext{}, &ChatGPTUserInfo{})
	if err != nil {
		common.Logln(2, "failed to migrate database")
		panic("failed to migrate database")
		return
	}
	DB = db
	var user []ChatGPTUserInfo
	result := DB.Find(&user)
	if result.Error != nil {
		common.Logln(2, "failed to find user info")
		panic("failed to find user info")
	}
	for _, v := range user {
		ChatGPTUsers[v.User] = v
	}
}
