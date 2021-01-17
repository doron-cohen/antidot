package utils

import (
	"os"
	"testing"
)

func TestPathExists(t *testing.T) {
	exists, err := PathExists("file_that_doesnt_exist")
	if err != nil {
		t.Fatalf("PathExists returned error: %v", err)
	}
	if exists {
		t.Fatalf("PathExists() returning true for file which doesn't exist")
	}

	f, err := os.Create("testfile")
	if err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}
	f.Close()
	defer os.Remove("testfile")

	exists, err = PathExists("testfile")
	if err != nil {
		t.Fatalf("PathExists returned error: %v", err)
	}
	if !exists {
		t.Fatalf("PathExists() returning false for file which exists")
	}

	os.RemoveAll("testdir") // Make sure it doesn't exist
	err = os.Mkdir("testdir", 0o755)
	if err != nil {
		t.Fatalf("Could not create test directory: %v", err)
	}
	defer os.RemoveAll("testdir")

	exists, err = PathExists("testdir")
	if err != nil {
		t.Fatalf("PathExists returned error: %v", err)
	}
	if exists {
		t.Fatalf("PathExists() returning true for empty directory")
	}

	f, err = os.Create("testdir/testfile")
	if err != nil {
		t.Fatalf("Error creating testdir/testfile")
	}
	f.Close()

	exists, err = PathExists("testdir")
	if err != nil {
		t.Fatalf("PathExists returned error: %v", err)
	}
	if !exists {
		t.Fatalf("PathExists() returning false for non-empty directory")
	}
}
