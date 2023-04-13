package common

import "os"

func GetFilesInfo(folderPath string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}
	return files, nil
}
