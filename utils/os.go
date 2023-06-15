package utils

import "os"

func CurrentUserHomePath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return dirname
}
