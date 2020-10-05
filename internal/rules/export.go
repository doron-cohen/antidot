package rules

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/doron-cohen/antidot/internal/tui"
	"github.com/doron-cohen/antidot/internal/utils"
)

type Export struct {
	Key   string
	Value string
}

func (e Export) Apply() error {
	envFile, err := utils.AppDirs.GetDataFile("env.sh")
	if err != nil {
		return err
	}

	if !utils.FileExists(envFile) {
		if _, err = os.Create(envFile); err != nil {
			return err
		}
	}

	envMap, err := godotenv.Read(envFile)
	if err != nil {
		return err
	}

	existingValue, isKeyContained := envMap[e.Key]
	if isKeyContained {
		log.Printf("Key %s already exists in env file %s", e.Key, envFile)
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

	log.Printf("Writing to %s", envFile)
	if err = utils.WriteEnvToFile(envMap, envFile); err != nil {
		return err
	}

	return nil
}

func (e Export) Pprint() {
	log.Printf(
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
