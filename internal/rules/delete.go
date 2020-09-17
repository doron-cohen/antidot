package rules

import (
	"log"
	"os"

	"github.com/doron-cohen/antidot/internal/utils"
)

type Delete struct {
	Path string
}

func (d Delete) Apply() error {
	path := os.ExpandEnv(d.Path)
	if !utils.FileExists(path) {
		return nil
	}

	err := os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func (d Delete) Pprint() {
	log.Printf("Delete %s", d.Path)
}
