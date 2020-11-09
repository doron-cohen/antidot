package rules

import (
	"errors"
	"fmt"
	"os"

	"github.com/doron-cohen/antidot/internal/tui"
	"github.com/doron-cohen/antidot/internal/utils"
)

type Export struct {
	Key   string
	Value string
}

func (e Export) Apply() error {
	envFile, err := utils.GetEnvFile()
	if err != nil {
		return err
	}

	if !utils.FileExists(envFile) {
		if _, err = os.Create(envFile); err != nil {
			return err
		}
	}

	envMap, err := utils.EnvMapFromFile(envFile)
	if err != nil {
		return err
	}

	existingValue, isKeyContained := envMap[e.Key]
	if isKeyContained {
		tui.Debug("Key %s already exists in env file %s", e.Key, envFile)
		if existingValue != e.Value {
			errMessage := fmt.Sprintf(
				"Current value for key '%s' (%s) is different than the requested (%s)",
				e.Key, existingValue, e.Value,
			)
			return errors.New(errMessage)
		} else {
			return nil
		}
	} else {
		envMap[e.Key] = e.Value
	}

	tui.Debug("Writing to %s", envFile)
	if err = utils.WriteKeyValuesToFile(envMap, envFile); err != nil {
		return err
	}

	return nil
}

func (e Export) Pprint() {
	tui.Print(
		"  %s %s%s\"%s\"",
		tui.ApplyStyle(tui.Blue, "EXPORT"),
		e.Key,
		tui.ApplyStyle(tui.Gray, "="),
		e.Value,
	)
}

func init() {
	registerAction("export", Export{})
}
