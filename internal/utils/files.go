package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/doron-cohen/antidot/internal/tui"
	"github.com/otiai10/copy"
)

func IsDir(filePath string) bool {
	info, err := os.Lstat(filePath)
	return err == nil && info.IsDir()
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func OpenOrCreate(filepath string) (*os.File, error) {
	if FileExists(filepath) {
		return os.Open(filepath)
	}
	return os.Create(filepath)
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

func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	if !fi.IsDir() {
		return true, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Treat empty directory as non-existent
	_, err = f.Readdirnames(1)
	if err != nil {
		if err == io.EOF {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Try to move file/directory with os.Rename and if that fails, do a copy + delete
func MovePath(source, dest string) error {
	exists, err := PathExists(dest)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Destination path %s already exists", dest)
	}

	err = os.Rename(source, dest)
	if err == nil {
		return nil
	}

	fi, err := os.Stat(source)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return MoveDirectory(source, dest)
	}

	return MoveFile(source, dest)
}

func AtomicWrite(bytes []byte, path string) error {
	tmpFilePath := fmt.Sprintf("%s.tmp", path)
	defer os.Remove(tmpFilePath)

	err := ioutil.WriteFile(tmpFilePath, bytes, os.FileMode(0o644))
	if err != nil {
		return err
	}

	err = os.Rename(tmpFilePath, path)
	if err != nil {
		return err
	}

	return nil
}
