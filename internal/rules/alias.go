package rules

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/doron-cohen/antidot/internal/tui"
	"github.com/doron-cohen/antidot/internal/utils"
)

type Alias struct {
	Alias   string
	Command string
}

// TODO: remove code duplication with export.go
func (a Alias) Apply() error {
	aliasFilePath, err := utils.GetAliasFile()
	if err != nil {
		return err
	}

	if !utils.FileExists(aliasFilePath) {
		if _, err = os.Create(aliasFilePath); err != nil {
			return err
		}
	}

	aliasMap, err := utils.AliasMapFromFile(aliasFilePath)
	if err != nil {
		return err
	}

	existingAlias, isAliasContained := aliasMap[a.Alias]
	if isAliasContained {
		log.Printf("Alias %s already exists in alias file %s", a.Alias, aliasMap)
		if existingAlias != a.Command {
			errMessage := fmt.Sprintf(
				"Current command for alias '%s' (%s) is different than the requested (%s)",
				a.Alias, existingAlias, a.Command,
			)
			return errors.New(errMessage)
		} else {
			return nil
		}
	} else {
		aliasMap[a.Alias] = a.Command
	}

	log.Printf("Writing to %s", aliasFilePath)
	if err = utils.WriteKeyValuesToFile(aliasMap, aliasFilePath); err != nil {
		return err
	}

	return nil
}

func (a Alias) Pprint() {
	log.Printf(
		"  %s %s%s\"%s\"",
		tui.ApplyStyle(tui.Magenta, "ALIAS"),
		a.Alias,
		tui.ApplyStyle(tui.Gray, "="),
		a.Command,
	)
}

func init() {
	registerAction("alias", Alias{})
}
