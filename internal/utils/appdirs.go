package utils

import (
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
