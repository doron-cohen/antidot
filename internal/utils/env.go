package utils

import (
	"fmt"
	"os"
	"os/user"

	"github.com/adrg/xdg"
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

func XdgVarsExport() string {
	return fmt.Sprintf(`export XDG_CONFIG_HOME="%s"
export XDG_CACHE_HOME="%s"
export XDG_DATA_HOME="%s"`,
		xdg.ConfigHome,
		xdg.CacheHome,
		xdg.DataHome,
	)
}

func GetAliasFile() (string, error) {
	return AppDirs.GetDataFile("alias.sh")
}

func GetEnvFile() (string, error) {
	return AppDirs.GetDataFile("env.sh")
}

func GetRulesFilePath() string {
	return AppDirs.GetDataFilePath("rules.yaml")
}
