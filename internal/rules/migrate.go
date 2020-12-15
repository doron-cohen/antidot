package rules

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/doron-cohen/antidot/internal/tui"
	"github.com/doron-cohen/antidot/internal/utils"
	"github.com/otiai10/copy"
)

type Migrate struct {
	Source  string
	Dest    string
	Symlink bool
}

func MoveDirectory(source, dest string) error {
	err := copy.Copy(source, dest)
	if err != nil {
		return err
	}

	err = os.RemoveAll(source)
	if err != nil {
		tui.Warn("Failed to remove original directory: %s", err)
	}

	return nil
}

// Adapted from https://gist.github.com/var23rav/23ae5d0d4d830aff886c3c970b8f6c6b
func MoveFile(source, dest string) error {
	inputFile, err := os.Open(source)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(dest)
	if err != nil {
		inputFile.Close()
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return err
	}

	// The copy was successful, so now delete the original file
	err = os.Remove(source)
	if err != nil {
		tui.Warn("Failed to remove original file: %s", err)
	}

	return nil
}

// Try to move file/directory with os.Rename and if that fails, do a copy + delete
func Move(source, dest string) error {
	err := os.Rename(source, dest)
	if err == nil {
		return nil
	}

	fi, err := os.Stat(source)
	if fi.Mode().IsDir() {
		return MoveDirectory(source, dest)
	}

	return MoveFile(source, dest)
}

func (m Migrate) Apply() error {
	source := utils.ExpandEnv(m.Source)
	_, err := os.Stat(source)
	if os.IsNotExist(err) {
		tui.Print("File %s doesn't exist. Skipping action", source)
		return nil
	} else if err != nil {
		return err
	}

	dest := utils.ExpandEnv(m.Dest)
	if utils.FileExists(dest) {
		errMessage := fmt.Sprintf("Destination file %s exists", dest)
		return errors.New(errMessage)
	}

	err = os.MkdirAll(filepath.Dir(dest), os.FileMode(0o755))
	if err != nil {
		return err
	}

	err = Move(source, dest)
	if err != nil {
		return err
	}

	if m.Symlink {
		err = os.Symlink(source, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m Migrate) Pprint() {
	symlink := ""
	if m.Symlink {
		symlink = " (with symlink)"
	}

	// TODO: move the indentation logic elsewhere
	tui.Print(
		"  %s %s %s %s%s",
		tui.ApplyStyle(tui.Green, "MOVE  "),
		utils.ExpandEnv(m.Source),
		tui.ApplyStyle(tui.Gray, "→"),
		utils.ExpandEnv(m.Dest),
		symlink)
}

func init() {
	registerAction("migrate", Migrate{})
}
