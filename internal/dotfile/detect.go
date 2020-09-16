package dotfile

import (
	"io/ioutil"
	"strings"
)

type Dotfile struct {
	Name  string
	IsDir bool
}

func NewDotfile(name string, isDir bool) *Dotfile {
	return &Dotfile{name, isDir}
}

func isDotfile(filename string) bool {
	return filename != "." && strings.HasPrefix(filename, ".")
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
		if isDotfile(filename) {
			dotfile := Dotfile{
				Name:  filename,
				IsDir: fileInfo.IsDir(),
			}
			found = append(found, dotfile)
		}
	}

	return found, nil
}
