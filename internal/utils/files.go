package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/doron-cohen/antidot/internal/tui"
	"github.com/otiai10/copy"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}

	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}

	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
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

// Try to move file/directory with os.Rename and if that fails, do a copy + delete
func MovePath(source, dest string) error {
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
