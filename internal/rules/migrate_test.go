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

func TestMigrateApply(t *testing.T) {
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

	sourcePath := filepath.Join(tmpDir, "source")
	_, err = os.Create(sourcePath)
	if err != nil {
		t.Errorf("Failed creating source file: %v", err)
	}

	destPath := filepath.Join(tmpDir, "dir/dest")
	migrateAction := rules.Migrate{Source: sourcePath, Dest: destPath, Symlink: false}
	err = migrateAction.Apply(&testActionContext)
	if err != nil {
		t.Fatalf("Error while applying migrate action: %v", err)
	}

	_, err = os.Stat(sourcePath)
	if os.IsExist(err) {
		t.Errorf("Source file %s still exists after Migrate action applied", sourcePath)
	}

	_, err = os.Stat(destPath)
	if os.IsNotExist(err) {
		t.Errorf("Destination file %s doesn't exist after Migrate action applied", sourcePath)
	}
}
