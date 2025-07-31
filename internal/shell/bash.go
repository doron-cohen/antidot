package shell

import (
	"fmt"
	"strings"

	"github.com/doron-cohen/antidot/internal/utils"
)

type Bash struct{}

func (b *Bash) formatAlias(alias, command string) string {
	return fmt.Sprintf("alias %s=\"%s\"\n", alias, command)
}

func (b *Bash) formatExport(key, value string) string {
	return fmt.Sprintf("export %s=\"%s\"\n", key, value)
}

func (b *Bash) RenderInit(kv *KeyValueStore) string {
	var sb strings.Builder

	// XDG exports
	for key, value := range utils.XdgDefaults() {
		sb.WriteString(fmt.Sprintf("export %s=\"${%s:-%s}\"\n", key, key, value))
	}

	if kv != nil {
		envs, _ := kv.ListEnvVars()
		for k, v := range envs {
			sb.WriteString(b.formatExport(k, v))
		}
		aliases, _ := kv.ListAliases()
		for k, v := range aliases {
			sb.WriteString(b.formatAlias(k, v))
		}
	}

	return sb.String()
}

// Compile-time interface check
var _ Shell = (*Bash)(nil)
