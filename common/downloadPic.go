package common

import (
	"io"
	"os"
)

// DownloadPic 下载网络图片到本地
func DownloadPic(savePath string, url string) error {
	resp, err := RequestTo(url, "GET", "", nil)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = os.WriteFile(savePath, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func DownloadPicWithProxy(savePath string, url string) error {
	resp, err := RequestTOWithProxy(url, "GET", "", nil)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = os.WriteFile(savePath, data, 0666)
	if err != nil {
		return err
	}
	return nil
}
