package rules

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/doron-cohen/antidot/internal/tui"
	"github.com/doron-cohen/antidot/internal/utils"
)

type Migrate struct {
	Source  string
	Dest    string
	Symlink bool
}

func (m Migrate) Apply() error {
	source := utils.ExpandEnv(m.Source)
	_, err := os.Stat(source)
	if os.IsNotExist(err) {
		log.Printf("File %s doesn't exist. Skipping action", source)
		return nil
	} else if err != nil {
		return err
	}

	dest := utils.ExpandEnv(m.Dest)
	if utils.FileExists(dest) {
		errMessage := fmt.Sprintf("Destination file %s exists", dest)
		return errors.New(errMessage)
	}

	err = os.MkdirAll(filepath.Dir(dest), os.FileMode(0o755))
	if err != nil {
		return err
	}

	err = os.Rename(source, dest)
	if err != nil {
		return err
	}

	if m.Symlink {
		err = os.Symlink(source, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m Migrate) Pprint() {
	symlink := ""
	if m.Symlink {
		symlink = " (with symlink)"
	}

	// TODO: move the indentation logic elsewhere
	log.Printf(
		"  %s %s %s %s%s",
		tui.ApplyStyle(tui.Green, "MOVE  "),
		utils.ExpandEnv(m.Source),
		tui.ApplyStyle(tui.Gray, "→"),
		utils.ExpandEnv(m.Dest),
		symlink)
}

func init() {
	registerAction("migrate", Migrate{})
}
