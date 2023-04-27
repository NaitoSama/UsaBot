package common

import (
	"encoding/base64"
	"io"
	"os"
	"strconv"
)

func At(tarID int64) string {
	return "[CQ:at,qq=" + strconv.FormatInt(tarID, 10) + "] "
}

func Pic(path string) string {
	return "[CQ:image,file=" + path + "] "
}

func PicBase64(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	picEncoded := base64.StdEncoding.EncodeToString(data)
	return Pic("base64://" + picEncoded), nil
}
