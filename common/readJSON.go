package common

import (
	"encoding/json"
	"os"
)

// ReadJSON 读取json文件并赋值给结构体
func ReadJSON(jsonPath string, bindStruct interface{}) error {
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, bindStruct)
	if err != nil {
		return err
	}
	return nil
}
