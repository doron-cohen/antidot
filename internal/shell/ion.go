package shell

import (
	"fmt"
	"strings"

	"github.com/doron-cohen/antidot/internal/utils"
)

type Ion struct{}

func (b *Ion) FormatAlias(alias, command string) string {
	return fmt.Sprintf("alias %s = \"%s\"\n", alias, command)
}

func (b *Ion) FormatExport(key, value string) string {
	// ion uses the variable $HISTFILE itself and it uses a proper location
	if (key == "HISTFILE") {
		return ""
	}
	return fmt.Sprintf("export %s = \"%s\"\n", key, value)
}

func (b *Ion) InitStub() string {
	builder := strings.Builder{}
	builder.WriteString("# Put 'eval $(antidot init)' (without single quotes) in your ionrc to automatically run this\n")
	for key, value := range utils.XdgDefaults() {
		builder.WriteString(fmt.Sprintf("export %s = $or(${%s} \"%s\")\n", key, key, value))
	}
	return builder.String()
}

func init() {
	registerShell("ion", &Ion{})
}
