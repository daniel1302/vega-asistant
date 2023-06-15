package utils

import "os"

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

func IsDir(filePath string) bool {
	stat, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return stat.IsDir()
}
