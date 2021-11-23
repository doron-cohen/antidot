package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/doron-cohen/antidot/internal/utils"
)

type Shell interface {
	FormatAlias(alias, command string) string
	FormatExport(alias, command string) string
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
		errorBuilder := strings.Builder{}
		errorBuilder.WriteString("Shell ")
		errorBuilder.WriteString(shellName)
		errorBuilder.WriteString("is not supported.\nSupported shells are:\n")
		for name := range SupportedShells {
			errorBuilder.WriteString(" ")
			errorBuilder.WriteString(name)
		}
		return nil, fmt.Errorf(errorBuilder.String())
	}

	return shell, nil
}

func GetShellScript(shell Shell) (string, error) {
	kvPath, err := utils.AppDirs.GetDataFile("kvstore.json")
	if err != nil {
		return "", err
	}

	kvStore, err := LoadKeyValueStore(kvPath)
	if err != nil {
		return "", err
	}
	if err := kvStore.load(); err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}
	builder.WriteString(shell.InitStub())
	builder.WriteString("\n")

	environment, err := DumpExports(shell, kvStore)
	if err != nil {
		return "", err
	}
	builder.WriteString(environment)
	builder.WriteString("\n")

	aliases, err := DumpAliases(shell, kvStore)
	if err != nil {
		return "", err
	}
	builder.WriteString(aliases)

	return builder.String(), nil
}

func ListShells() string {
	builder := strings.Builder{}
	for name, _ := range SupportedShells {
		builder.WriteString(" ")
		builder.WriteString(name)
	}
	return builder.String()
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

func DumpAliases(shell Shell, kvStore *KeyValueStore) (string, error) {
	aliases, err := kvStore.ListAliases()
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}
	for k, v := range aliases {
		kvLine := shell.FormatAlias(k, v)
		builder.WriteString(kvLine)
	}

	return builder.String(), nil
}

func DumpExports(shell Shell, kvStore *KeyValueStore) (string, error) {
	envVars, err := kvStore.ListEnvVars()
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}
	for k, v := range envVars {
		kvLine := shell.FormatExport(k, v)
		builder.WriteString(kvLine)
	}

	return builder.String(), nil
}
