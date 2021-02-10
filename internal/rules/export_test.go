package rules_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/doron-cohen/antidot/internal/rules"
	"github.com/doron-cohen/antidot/internal/shell"
	"github.com/doron-cohen/antidot/internal/utils"
)

func TestExportNoEnv(t *testing.T) {
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
	exportAction := rules.Export{Key: "HELLO", Value: "world"}
	err = exportAction.Apply(&testActionContext)
	if err != nil {
		t.Fatalf("Error while applying export action: %v", err)
	}

	envVars, err := kvStore.ListEnvVars()
	if err != nil {
		t.Fatalf("Error while listing env vars from kv store: %v", err)
	}

	if !cmp.Equal(envVars, map[string]string{"HELLO": "world"}) {
		t.Fatalf("Unexpected env file contents: %s", envVars)
	}
}

// TODO: test env key conflict
