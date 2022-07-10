package dotfile_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/doron-cohen/antidot/internal/dotfile"
)

type fileInfo struct {
	path  string
	isDir bool
}

type dotfilesTest struct {
	Name             string
	Files            []fileInfo
	ExpectedDotfiles []dotfile.Dotfile
}

func setup(files []fileInfo) (string, error) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		return "", err
	}

	for _, f := range files {
		path := filepath.Join(tmpDir, f.path)

		var dir string
		if f.isDir {
			dir = path
		} else {
			dir = filepath.Dir(path)
		}

		err := os.MkdirAll(dir, os.FileMode(0o755))
		if err != nil {
			return "", err
		}

		if !f.isDir {
			_, err := os.Create(path)
			if err != nil {
				return "", err
			}
		}
	}
	return tmpDir, nil
}

func tearDown(tmpDir string) {
	if err := os.RemoveAll(tmpDir); err != nil {
		log.Printf("Failed to cleanup temp dir %s", tmpDir)
	}
}

func TestFindDotfiles(t *testing.T) {
	tests := []dotfilesTest{
		{
			Name:             "Empty dir",
			Files:            []fileInfo{},
			ExpectedDotfiles: []dotfile.Dotfile{},
		},
		{
			Name: "No dotfiles",
			Files: []fileInfo{
				{"file1", false},
				{"file2", false},
				{"dir1", true},
				{"dir1/file3", false},
				{"dir1/file4", false},
			},
			ExpectedDotfiles: []dotfile.Dotfile{
				{"dir1", true},
				{"file1", false},
				{"file2", false},
			},
		},
		{
			Name: "Mixed",
			Files: []fileInfo{
				{"file1", false},
				{".file2", false},
				{"dir1", true},
				{"dir1/.file3", false},
				{"dir1/file4", false},
				{".dir2", true},
				{".dir2/.file5", true},
			},
			ExpectedDotfiles: []dotfile.Dotfile{
				// Order matters here and I don't want to invest in a comparison function
				{".dir2", true},
				{".file2", false},
				{"dir1", true},
				{"file1", false},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			tmpDir, err := setup(test.Files)
			defer tearDown(tmpDir)

			if err != nil {
				t.Errorf("Failed setting up test: %v", err)
			}

			found, err := dotfile.Detect(tmpDir)
			if err != nil {
				t.Errorf("Failed to detect dotfiles from %s: %v", tmpDir, err)
			}

			if !cmp.Equal(found, test.ExpectedDotfiles) {
				t.Errorf("Found dotfiles are not what we expected:\n%v", cmp.Diff(found, test.ExpectedDotfiles))
			}
		})
	}
}
