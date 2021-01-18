package shell_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/doron-cohen/antidot/internal/shell"
	"github.com/google/go-cmp/cmp"
)

func appendToEndOfFile(path, text string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.FileMode(0o755))
	defer func() { _ = file.Close() }()
	if err != nil {
		return fmt.Errorf("Failed to open key value file: %v", err)
	}

	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("Failed to seek key value file: %v", err)
	}

	_, err = file.WriteString(text)
	if err != nil {
		return fmt.Errorf("Failed to write to key value file: %v", err)
	}

	return nil
}

type kvTest struct {
	testName          string
	key               string
	value             string
	expected          map[string]string
	expectedAppendErr string
	appendText        string
}

func TestKeyValueFile(t *testing.T) {
	expected := make(map[string]string)
	format := "set %s = '%s'\n"
	kv := shell.NewKeyValueMapFormat(
		`(?m)^set (?P<key>\w+) = '(?P<value>.*)'`,
		format,
	)

	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Errorf("Failed setting up test: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	path := filepath.Join(tmpDir, "test.kv")

	tests := []kvTest{
		{
			testName: "add key value to empty file",
			key:      "hello",
			value:    "world",
			expected: map[string]string{"hello": "world"},
		},
		{
			testName: "add key value to existing file",
			key:      "hi",
			value:    "there",
			expected: map[string]string{"hello": "world", "hi": "there"},
		},
		{
			testName:   "ignore out of format text",
			key:        "",
			value:      "",
			expected:   map[string]string{"hello": "world", "hi": "there"},
			appendText: "This is not the format!\n",
		},
		{
			testName:          "key value conflict",
			key:               "hi",
			value:             "ho",
			expected:          map[string]string{"hello": "world", "hi": "there"},
			expectedAppendErr: "Key hi already exists with different value",
		},
		{
			testName: "ignore key value exists",
			key:      "hi",
			value:    "there",
			expected: map[string]string{"hello": "world", "hi": "there"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if tt.appendText != "" {
				err = appendToEndOfFile(path, tt.appendText)
				if err != nil {
					t.Fatal(err)
				}
			}

			if tt.key != "" {
				err = shell.AppendKeyValueToFile(path, tt.key, tt.value, kv)

				if tt.expectedAppendErr != "" {
					if err == nil || err.Error() != tt.expectedAppendErr {
						t.Fatalf("Expected error %s from append key value", tt.expectedAppendErr)
					}
					return
				}
				if err != nil {
					t.Fatal("Failed to append key value to file", err)
				}
			}

			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatal("Failed to read key value file")
			}
			kvMap, err := kv.ParseKeyValuePairs(string(bytes[:]))
			if err != nil {
				t.Fatal("Failed to parse key values from file", err)
			}
			if !cmp.Equal(kvMap, tt.expected) {
				t.Fatalf("Key value mismatch. Expected %v got %v", expected, kvMap)
			}
		})
	}
}
