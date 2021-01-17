package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Shell interface {
	AddEnvExport(key, value string) error
	AddCommandAlias(alias, command string) error
	InitStub() string
}

var SupportedShells = make(map[string]Shell)
var FallbackShellName = "bash"

func Get(shellName string) (Shell, error) {
	if shellName == "" {
		shellName = detectShell()
		if shellName == "" {
			shellName = FallbackShellName
		}
	}

	shell, ok := SupportedShells[shellName]
	if !ok {
		return nil, fmt.Errorf("Shell %s is still not supported.", shellName)
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

func registerShell(name string, shell Shell) {
	SupportedShells[name] = shell
}
