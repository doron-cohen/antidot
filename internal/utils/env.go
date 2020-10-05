package utils

import (
	"fmt"
	"os"
	"os/user"

	"github.com/joho/godotenv"
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

func WriteEnvToFile(envMap map[string]string, filePath string) error {
	var newKey string
	newMap := make(map[string]string, len(envMap))
	for key, value := range envMap {
		// TODO: remove this ugly but working hack
		newKey = fmt.Sprintf("export %s", key)
		newMap[newKey] = value
	}

	return godotenv.Write(newMap, filePath)
}
