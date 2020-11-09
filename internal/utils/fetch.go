package utils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/doron-cohen/antidot/internal/tui"
)

func Download(src, dest string) error {
	resp, err := http.Get(src)
	if err != nil {
		return err
	}
	// TODO: check for close() errors (across the code)
	defer resp.Body.Close()

	tempFile, err := ioutil.TempFile("", "rules.*.yaml")
	tui.FatalIfError("Failed to create rules file", err)
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return err
	}

	dir := path.Dir(dest)
	fileInfo, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, os.FileMode(0o755)); err != nil {
			return err
		}
	} else if !fileInfo.IsDir() {
		text := fmt.Sprintf("Rules file destination directory is a file: %s", dir)
		return errors.New(text)
	}

	err = MoveFile(tempFile.Name(), dest)
	if err != nil {
		return err
	}

	return nil
}
