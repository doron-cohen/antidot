package rules

import (
	"log"
	"os"

	"github.com/doron-cohen/antidot/internal/tui"
	"github.com/doron-cohen/antidot/internal/utils"
)

type Delete struct {
	Path string
}

func (d Delete) Apply() error {
	path := utils.ExpandEnv(d.Path)
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
	log.Printf("  %s %s", tui.ApplyStyle(tui.Red, "DELETE"), utils.ExpandEnv(d.Path))
}
