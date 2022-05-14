package dotfile

import (
	"io/ioutil"
)

type Dotfile struct {
	Name  string
	IsDir bool `mapstructure:"is_dir"`
}

func Detect(dir string) ([]Dotfile, error) {
	// TODO: better handle file errors
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	found := make([]Dotfile, 0, len(files))
	for _, fileInfo := range files {
		filename := fileInfo.Name()
		dotfile := Dotfile{
			Name:  filename,
			IsDir: fileInfo.IsDir(),
		}
		found = append(found, dotfile)
	}

	return found, nil
}
