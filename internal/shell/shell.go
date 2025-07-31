package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Shell interface {
	RenderInit(kv *KeyValueStore) string
}

var FallbackShellName = "bash"

func Get(shellName string) (Shell, error) {
	if shellName == "" {
		shellName = detectShell()
		if shellName == "" {
			shellName = FallbackShellName
		}
	}

	var shell Shell
	switch shellName {
	case "bash":
		shell = &Bash{}
	case "fish":
		shell = &Fish{}
	case "zsh":
		shell = &Bash{} // Zsh uses the same implementation as bash
	default:
		return nil, fmt.Errorf("Shell %s is not supported. Supported shells: [bash fish zsh]", shellName)
	}

	return shell, nil
}

func detectShell() string {
	shellPath := os.Getenv("SHELL")
	if shellPath == "" {
		return ""
	}

	return filepath.Base(strings.Split(shellPath, " ")[0])
}
