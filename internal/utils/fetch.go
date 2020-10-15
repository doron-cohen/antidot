package utils

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func Download(src, dest string) error {
	resp, err := http.Get(src)
	if err != nil {
		return err
	}
	// TODO: check for close() errors (across the code)
	defer resp.Body.Close()

	tempFile, err := ioutil.TempFile("", "rules.*.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return err
	}

	err = MoveFile(tempFile.Name(), dest)
	if err != nil {
		return err
	}

	return nil
}
