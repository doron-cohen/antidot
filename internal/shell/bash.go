package shell

import (
	"fmt"

	"github.com/doron-cohen/antidot/internal/utils"
)

type Bash struct {
	envMapFormat   *keyValueMapFormat
	aliasMapFormat *keyValueMapFormat
}

func (b *Bash) envFilePath() (string, error) {
	return utils.AppDirs.GetDataFile("env.sh")
}

func (b *Bash) aliasFilePath() (string, error) {
	return utils.AppDirs.GetDataFile("alias.sh")
}

func (b *Bash) AddEnvExport(key, value string) error {
	path, err := b.envFilePath()
	if err != nil {
		return err
	}

	return AppendKeyValueToFile(path, key, value, b.envMapFormat)
}

func (b *Bash) AddCommandAlias(alias, command string) error {
	path, err := b.aliasFilePath()
	if err != nil {
		return err
	}

	return AppendKeyValueToFile(path, alias, command, b.aliasMapFormat)
}

func (b *Bash) InitStub() string {
	envFilePath, _ := b.envFilePath()
	aliasFilePath, _ := b.aliasFilePath()

	xdgExport := ""
	for key, value := range utils.XdgDefaults() {
		xdgExport += fmt.Sprintf("export %s=\"{%s:-%s}\"\n", key, key, value)
	}

	return fmt.Sprintf(`%s
if [ -f "%s" ]; then source "%s"; fi
if [ -f "%s" ]; then source "%s"; fi`,
		xdgExport,
		envFilePath,
		envFilePath,
		aliasFilePath,
		aliasFilePath,
	)
}

func init() {
	registerShell("bash", &Bash{
		envMapFormat: NewKeyValueMapFormat(
			`^export (?P<key>\w+)="(?P<value>.*)"`,
			"export %s=\"%s\"\n",
		),
		aliasMapFormat: NewKeyValueMapFormat(
			`^alias (?P<alias>\w+)="(?P<command>.*)"`,
			"alias %s=\"%s\"\n",
		),
	})
}
