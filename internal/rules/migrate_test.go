package rules_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/doron-cohen/antidot/internal/rules"
)

func TestMigrateApply(t *testing.T) {
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
	migrateAction.Apply()

	_, err = os.Stat(sourcePath)
	if os.IsExist(err) {
		t.Errorf("Source file %s still exists after Migrate action applied", sourcePath)
	}

	_, err = os.Stat(destPath)
	if os.IsNotExist(err) {
		t.Errorf("Destination file %s doesn't exist after Migrate action applied", sourcePath)
	}
}
