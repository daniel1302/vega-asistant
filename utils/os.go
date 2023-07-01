package utils

import (
	"fmt"
	"os"
	"os/user"
	"strings"
)

func CurrentUserHomePath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return dirname
}

func Whoami() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get the current user: %w", err)
	}

	username := currentUser.Username
	return username, nil
}

func IsWSL() bool {
	versionContent, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}

	return strings.Contains(string(versionContent), "microsoft") ||
		strings.Contains(string(versionContent), "WSL")
}
