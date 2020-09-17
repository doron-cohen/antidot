package action

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/doron-cohen/antidot/internal/utils.go"

	"github.com/mitchellh/mapstructure"
)

type Action interface {
	Apply() error
	Pprint()
}

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

func getActionByName(name string) (Action, error) {
	switch name {
	case "migrate":
		return Migrate{}, nil
	default:
		errMessage := fmt.Sprintf("Unknown action type '%s'", name)
		return nil, errors.New(errMessage)
	}
}

func actionDecodeHook(sourceType, destType reflect.Type, raw interface{}) (interface{}, error) {
	// TODO: find a better way to compare these types
	if fmt.Sprintf("%s", destType) == "action.Action" {
		var err error
		var result Action

		rawMap := raw.(map[interface{}]interface{})
		result, err = getActionByName(rawMap["type"].(string))
		if err != nil {
			return nil, err
		}

		mapstructure.Decode(raw, &result)
		return result, nil
	}
	return raw, nil
}
