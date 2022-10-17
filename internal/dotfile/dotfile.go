package dotfile

import (
	"os"
)

type Dotfile struct {
	Name  string
	IsDir bool `mapstructure:"is_dir"`
}

func (d *Dotfile) MatchPath(filePath string) (bool, error) {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir() == d.IsDir, nil
}
