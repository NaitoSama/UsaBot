package common

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ModifyPicMD5 修改图片的MD5
func ModifyPicMD5(picPath string) (string, error) {
	data, err := os.ReadFile(picPath)
	if err != nil {
		return "", err
	}
	modifiedData := append(data, []byte(time.Now().String())...)
	fileDir := filepath.Dir(picPath)
	fileName := filepath.Base(picPath)
	fileExtension := filepath.Ext(picPath)
	newFileName := fmt.Sprintf("%smodifiedMD5%s", fileName[:len(fileName)-len(fileExtension)], fileExtension)
	picNewPath := filepath.Join(fileDir, newFileName)
	err = os.WriteFile(picNewPath, modifiedData, os.ModePerm)
	if err != nil {
		return "", err
	}
	return picNewPath, nil
}
