package shell_test

import (
	"os"
	"testing"

	"github.com/doron-cohen/antidot/internal/shell"
)

type TestSh struct{}

func (t *TestSh) EnvFilePath() (string, error) {
	return "", nil
}

func (t *TestSh) AliasFilePath() (string, error) {
	return "", nil
}

func (t *TestSh) FormatAlias(alias, command string) string {
	return ""
}

func (t *TestSh) FormatExport(key, value string) string {
	return ""
}

func (t *TestSh) InitStub() string {
	return ""
}

type FallbackSh struct{}

func (t *FallbackSh) EnvFilePath() (string, error) {
	return "", nil
}

func (t *FallbackSh) AliasFilePath() (string, error) {
	return "", nil
}

func (t *FallbackSh) FormatAlias(alias, command string) string {
	return ""
}

func (t *FallbackSh) FormatExport(key, value string) string {
	return ""
}
func (t *FallbackSh) InitStub() string {
	return ""
}

func TestGetShell(t *testing.T) {
	testsh := TestSh{}
	fallbackSh := FallbackSh{}
	shell.SupportedShells = make(map[string]shell.Shell, 1)
	shell.SupportedShells["testsh"] = &testsh
	shell.SupportedShells[shell.FallbackShellName] = &fallbackSh

	_, err := shell.Get("unknownsh")
	if err == nil {
		t.Fatal("Expected error for unknown shell")
	}

	sh, err := shell.Get("testsh")
	if err != nil {
		t.Fatal(err)
	}
	if sh != &testsh {
		t.Fatalf("Unexpected shell. Expected %v got %v", testsh, sh)
	}

	os.Setenv("SHELL", "testsh")
	sh, err = shell.Get("")
	if err != nil {
		t.Fatal(err)
	}
	if sh != &testsh {
		t.Fatalf("Unexpected shell. Expected %v got %v", testsh, sh)
	}

	os.Unsetenv("SHELL")
	sh, err = shell.Get("")
	if err != nil {
		t.Fatal(err)
	}
	if sh != &fallbackSh {
		t.Fatalf("Unexpected shell. Expected %v got %v", fallbackSh, sh)
	}
}
