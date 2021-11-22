package shell

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/doron-cohen/antidot/internal/tui"
)

type Shell interface {
	AliasFilePath() (string, error)
	EnvFilePath() (string, error)
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
		var error = strings.Join([]string{"Shell", shellName, "is not supported.\nSupported shells are:\n",}, " ")
		for name := range SupportedShells {
			error = strings.Join([]string{error, name}, " ");
		}
		return nil, fmt.Errorf(error)
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

func DumpAliases(shell Shell, kvStore *KeyValueStore) error {
	aliasFilePath, err := shell.AliasFilePath()
	if err != nil {
		return err
	}

	aliases, err := kvStore.ListAliases()
	if err != nil {
		return err
	}

	builder := strings.Builder{}
	for k, v := range aliases {
		kvLine := shell.FormatAlias(k, v)
		builder.WriteString(kvLine)
	}

	tui.Debug("Dumping aliases to %s", aliasFilePath)
	err = ioutil.WriteFile(aliasFilePath, []byte(builder.String()), os.FileMode(0o644))
	if err != nil {
		return err
	}

	return nil
}

func DumpExports(shell Shell, kvStore *KeyValueStore) error {
	envFilePath, err := shell.EnvFilePath()
	if err != nil {
		return err
	}

	envVars, err := kvStore.ListEnvVars()
	if err != nil {
		return err
	}

	builder := strings.Builder{}
	for k, v := range envVars {
		kvLine := shell.FormatExport(k, v)
		builder.WriteString(kvLine)
	}

	tui.Debug("Dumping env vars to %s", envFilePath)
	err = ioutil.WriteFile(envFilePath, []byte(builder.String()), os.FileMode(0o644))
	if err != nil {
		return err
	}

	return nil
}
