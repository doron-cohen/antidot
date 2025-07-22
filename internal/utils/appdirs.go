package utils

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

type appDirs struct {
	AppName string
}

func (a appDirs) ConfigHome() string {
	return filepath.Join(xdg.ConfigHome, a.AppName)
}

func (a appDirs) CacheHome() string {
	return filepath.Join(xdg.CacheHome, a.AppName)
}

func (a appDirs) DataHome() string {
	return filepath.Join(xdg.DataHome, a.AppName)
}

func (a appDirs) GetDataFilePath(fileName string) string {
	return filepath.Join(a.DataHome(), fileName)
}

func (a appDirs) GetDataFile(fileName string) (string, error) {
	relPath := filepath.Join(a.AppName, fileName)
	return xdg.DataFile(relPath)
}

var AppDirs appDirs

func init() {
	// TODO: import this from somewhere
	AppDirs = appDirs{AppName: "antidot"}
}

func GetKeyValueStorePath() (string, error) {
	// Check if ANTIDOT_STATE_FILE environment variable is set
	if stateFile := os.Getenv("ANTIDOT_STATE_FILE"); stateFile != "" {
		return stateFile, nil
	}

	// Fall back to default XDG location
	return AppDirs.GetDataFile("kvstore.json")
}
