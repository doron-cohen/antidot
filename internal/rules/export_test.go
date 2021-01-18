package rules_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/doron-cohen/antidot/internal/rules"
	"github.com/doron-cohen/antidot/internal/utils"
)

func TestExportNoEnv(t *testing.T) {
	utils.AppDirs.AppName = "antidot_test"
	defer os.RemoveAll(utils.AppDirs.DataHome())

	envFilePath, err := utils.GetEnvFile()
	if err != nil {
		t.Fatalf("Error getting env file path: %v", err)
	}

	if utils.FileExists(envFilePath) {
		t.Fatalf("Env file %s shouldn't exist", envFilePath)
	}

	exportAction := rules.Export{Key: "HELLO", Value: "world"}
	err = exportAction.Apply(testActionContext)
	if err != nil {
		t.Fatalf("Error while applying export action: %v", err)
	}

	contents, err := ioutil.ReadFile(envFilePath)
	if err != nil {
		t.Fatalf("Error while reading env file: %v", err)
	}

	// TODO: don't be too specific
	if !cmp.Equal(contents, []byte("export HELLO=\"world\"\n")) {
		t.Fatalf("Unexpected env file contents: %s", contents)
	}
}

// TODO: test env key conflict
