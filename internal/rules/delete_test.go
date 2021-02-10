package rules_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/doron-cohen/antidot/internal/rules"
	"github.com/doron-cohen/antidot/internal/shell"
	"github.com/doron-cohen/antidot/internal/utils"
)

func TestDeleteApply(t *testing.T) {
	utils.AppDirs.AppName = "antidot_test"
	defer os.RemoveAll(utils.AppDirs.DataHome())

	kvPath := filepath.Join(utils.AppDirs.DataHome(), "store.json")
	kvStore, _ := shell.LoadKeyValueStore(kvPath)
	testActionContext := rules.ActionContext{kvStore}

	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Errorf("Failed setting up test: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "file")
	_, err = os.Create(filePath)
	if err != nil {
		t.Errorf("Failed creating test file: %v", err)
	}

	deleteAction := rules.Delete{Path: filePath}
	err = deleteAction.Apply(&testActionContext)
	if err != nil {
		t.Fatalf("Error while applying delete action: %v", err)
	}

	_, err = os.Stat(filePath)
	if os.IsExist(err) {
		t.Errorf("File %s still exists after Delete action applied", filePath)
	}
}
