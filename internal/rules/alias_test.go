package rules_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/doron-cohen/antidot/internal/rules"
	"github.com/doron-cohen/antidot/internal/utils"
)

func TestAliasNoEnv(t *testing.T) {
	utils.AppDirs.AppName = "antidot_test"
	defer os.RemoveAll(utils.AppDirs.DataHome())

	aliasFilePath, err := utils.GetAliasFile()
	if err != nil {
		t.Fatalf("Error getting alias file path: %v", err)
	}

	if utils.FileExists(aliasFilePath) {
		t.Fatalf("Alias file %s shouldn't exist", aliasFilePath)
	}

	aliasAction := rules.Alias{Alias: "ll", Command: "ls -la \"$XDG_CONFIG_HOME\"/test"}
	err = aliasAction.Apply()
	if err != nil {
		t.Fatalf("Error while applying alias action: %v", err)
	}

	contents, err := ioutil.ReadFile(aliasFilePath)
	if err != nil {
		t.Fatalf("Error while reading alias file: %v", err)
	}

	// TODO: don't be too specific
	expected := "alias ll=\"ls -la \"$XDG_CONFIG_HOME\"/test\""
	if !cmp.Equal(strings.Trim(string(contents), " \t\n"), string(expected)) {
		t.Fatalf("Unexpected alias file contents: '%s' != '%s'", contents, expected)
	}
}

// TODO: test alias key conflict
