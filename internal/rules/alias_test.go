package rules_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/doron-cohen/antidot/internal/rules"
	"github.com/doron-cohen/antidot/internal/shell"
	"github.com/doron-cohen/antidot/internal/utils"
)

func TestAliasNoEnv(t *testing.T) {
	utils.AppDirs.AppName = "antidot_test"

	kvPath, err := utils.AppDirs.GetDataFile("store.json")
	if err != nil {
		t.Fatalf("Failed getting data file path %s: %s", kvPath, err)
	}
	defer os.RemoveAll(utils.AppDirs.DataHome())

	kvStore, err := shell.LoadKeyValueStore(kvPath)
	if err != nil {
		t.Fatalf("Failed loading key value store from %s: %s", kvPath, err)
	}

	testActionContext := rules.ActionContext{kvStore}
	aliasAction := rules.Alias{Alias: "ll", Command: "ls -la \"$XDG_CONFIG_HOME\"/test"}
	err = aliasAction.Apply(&testActionContext)
	if err != nil {
		t.Fatalf("Error while applying alias action: %v", err)
	}

	aliases, err := kvStore.ListAliases()
	if err != nil {
		t.Fatalf("Error while listing env vars from kv store: %v", err)
	}

	if !cmp.Equal(aliases, map[string]string{"ll": "ls -la \"$XDG_CONFIG_HOME\"/test"}) {
		t.Fatalf("Unexpected alias file contents: %s", aliases)
	}
}

// TODO: test alias key conflict
