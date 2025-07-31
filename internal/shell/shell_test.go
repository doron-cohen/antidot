package shell_test

import (
	"os"
	"testing"

	"github.com/doron-cohen/antidot/internal/shell"
)

type TestSh struct{}

func (t *TestSh) RenderInit(kv *shell.KeyValueStore) string {
	return ""
}

type FallbackSh struct{}

func (t *FallbackSh) RenderInit(kv *shell.KeyValueStore) string {
	return ""
}

func TestGetShell(t *testing.T) {
	// Test that we can get all supported shells
	bash, err := shell.Get("bash")
	if err != nil {
		t.Fatal(err)
	}
	if bash == nil {
		t.Fatal("Expected bash shell, got nil")
	}

	fish, err := shell.Get("fish")
	if err != nil {
		t.Fatal(err)
	}
	if fish == nil {
		t.Fatal("Expected fish shell, got nil")
	}

	zsh, err := shell.Get("zsh")
	if err != nil {
		t.Fatal(err)
	}
	if zsh == nil {
		t.Fatal("Expected zsh shell, got nil")
	}

	// Test that zsh returns the same as bash (since they share implementation)
	if bash != zsh {
		t.Fatalf("Expected zsh and bash to be the same type, got different types")
	}

	// Test unknown shell
	_, err = shell.Get("unknownsh")
	if err == nil {
		t.Fatal("Expected error for unknown shell")
	}

	// Test shell detection from environment
	os.Setenv("SHELL", "bash")
	detected, err := shell.Get("")
	if err != nil {
		t.Fatal(err)
	}
	if detected == nil {
		t.Fatal("Expected detected shell, got nil")
	}

	// Test fallback when no shell is detected
	os.Unsetenv("SHELL")
	fallback, err := shell.Get("")
	if err != nil {
		t.Fatal(err)
	}
	if fallback == nil {
		t.Fatal("Expected fallback shell, got nil")
	}
}
