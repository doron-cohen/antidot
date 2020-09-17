package rules

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/doron-cohen/antidot/internal/utils"
)

type Migrate struct {
	Source  string
	Dest    string
	Symlink bool
}

func (m Migrate) Apply() error {
	source := os.ExpandEnv(m.Source)
	_, err := os.Stat(source)
	if os.IsNotExist(err) {
		log.Printf("File %s doesn't exist. Skipping action", source)
		return nil
	} else if err != nil {
		return err
	}

	dest := os.ExpandEnv(m.Dest)
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
	log.Printf("Move %s to %s. Symlink: %v", m.Source, m.Dest, m.Symlink)
}
