package utils

import (
	"os"
	"os/user"
)

func GetHomeDir() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	return user.HomeDir, nil
}

func ExpandEnv(text string) string {
	return os.ExpandEnv(text)
}
