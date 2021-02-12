package shell

import (
	"fmt"

	"github.com/doron-cohen/antidot/internal/utils"
)

type Bash struct{}

func (b *Bash) EnvFilePath() (string, error) {
	return utils.AppDirs.GetDataFile("env.sh")
}

func (b *Bash) AliasFilePath() (string, error) {
	return utils.AppDirs.GetDataFile("alias.sh")
}

func (b *Bash) FormatAlias(alias, command string) string {
	return fmt.Sprintf("alias %s=\"%s\"\n", alias, command)
}

func (b *Bash) FormatExport(key, value string) string {
	return fmt.Sprintf("export %s=\"%s\"\n", key, value)
}

func (b *Bash) InitStub() string {
	envFilePath, _ := b.EnvFilePath()
	aliasFilePath, _ := b.AliasFilePath()

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
	registerShell("bash", &Bash{})
}
