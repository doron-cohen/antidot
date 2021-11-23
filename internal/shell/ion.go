package shell

import (
	"fmt"
	"strings"

	"github.com/doron-cohen/antidot/internal/utils"
)

type Bash struct{}

func (b *Bash) FormatAlias(alias, command string) string {
	return fmt.Sprintf("alias %s = \"%s\"\n", alias, command)
}

func (b *Bash) FormatExport(key, value string) string {
	return fmt.Sprintf("export %s = \"%s\"\n", key, value)
}

func (b *Bash) InitStub() string {
	builder := strings.Builder{}
	builder.WriteString("# Put 'eval \"$(antidot init -c ion)\"' (without single quotes) in your ionrc to automatically run this\n")
	for key, value := range utils.XdgDefaults() {
		builder.WriteString(fmt.Sprintf("export %s = $or(${%s} \"%s\")\n", key, key, value))
	}
	return builder.String()
}

func init() {
	registerShell("bash", &Bash{})
}
