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
		dotfilesTest{
			Name:             "Empty dir",
			Files:            []fileInfo{},
			ExpectedDotfiles: []dotfile.Dotfile{},
		},
		dotfilesTest{
			Name: "No dotfiles",
			Files: []fileInfo{
				fileInfo{"file1", false},
				fileInfo{"file2", false},
				fileInfo{"dir1", true},
				fileInfo{"dir1/file3", false},
				fileInfo{"dir1/file4", false},
			},
			ExpectedDotfiles: []dotfile.Dotfile{},
		},
		dotfilesTest{
			Name: "Mixed",
			Files: []fileInfo{
				fileInfo{"file1", false},
				fileInfo{".file2", false},
				fileInfo{"dir1", true},
				fileInfo{"dir1/.file3", false},
				fileInfo{"dir1/file4", false},
				fileInfo{".dir2", true},
				fileInfo{".dir2/.file5", true},
			},
			ExpectedDotfiles: []dotfile.Dotfile{
				// Order matters here and I don't want to invest in a comparison function
				dotfile.Dotfile{".dir2", true},
				dotfile.Dotfile{".file2", false},
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
