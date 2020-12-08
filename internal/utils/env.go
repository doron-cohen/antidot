package utils

import (
	"fmt"
	"os"
	"os/user"

	"github.com/adrg/xdg"

	"github.com/doron-cohen/antidot/internal/tui"
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

func ApplyDefaultXdgEnv() {
	xdgSystemDefaults := map[string]string{
		"XDG_CONFIG_HOME": xdg.ConfigHome,
		"XDG_CACHE_HOME":  xdg.CacheHome,
		"XDG_DATA_HOME":   xdg.DataHome,
	}
	printNewline := false
	for name, defaultValue := range xdgSystemDefaults {
		if value, exists := os.LookupEnv(name); !exists || value == "" {
			tui.Warn(
				"Environment variable %s not set. Using default path: %s",
				tui.ApplyStyle(tui.Yellow, name),
				tui.ApplyStyle(tui.Yellow, defaultValue),
			)
			os.Setenv(name, defaultValue)
			printNewline = true
		}
	}

	if printNewline {
		fmt.Println("")
	}
}

func XdgVarsExport() string {
	return fmt.Sprintf(`export XDG_CONFIG_HOME="${XDG_CONFIG_HOME:-"%s"}"
export XDG_CACHE_HOME="${XDG_CACHE_HOME:-"%s"}"
export XDG_DATA_HOME="${XDG_DATA_HOME:-"%s"}"`,
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
