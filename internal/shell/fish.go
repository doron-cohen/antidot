package shell

import (
	"fmt"
	"strings"

	"github.com/doron-cohen/antidot/internal/utils"
)

type Fish struct{}

func (f *Fish) formatAlias(alias, command string) string {
	unbracketed := f.unbracketEnvVar(command)
	return fmt.Sprintf("alias %s \"%s\"\n", alias, unbracketed)
}

func (f *Fish) formatExport(key, value string) string {
	unbracketed := f.unbracketEnvVar(value)
	return fmt.Sprintf("set -gx %s \"%s\"\n", key, unbracketed)
}

func (f *Fish) unbracketEnvVar(value string) string {
	// Replace ${VAR} with $VAR for fish shell
	// First replace ${ with $
	result := strings.ReplaceAll(value, "${", "$")
	// Then remove the closing }
	result = strings.ReplaceAll(result, "}", "")
	return result
}

func (f *Fish) RenderInit(kv *KeyValueStore) string {
	var sb strings.Builder

	// XDG exports with fish syntax
	for key, value := range utils.XdgDefaults() {
		sb.WriteString(fmt.Sprintf("set -q %s; or set -x %s \"%s\"\n", key, key, value))
	}

	if kv != nil {
		envs, _ := kv.ListEnvVars()
		for k, v := range envs {
			sb.WriteString(f.formatExport(k, v))
		}
		aliases, _ := kv.ListAliases()
		for k, v := range aliases {
			sb.WriteString(f.formatAlias(k, v))
		}
	}

	return sb.String()
}

// Compile-time interface check
var _ Shell = (*Fish)(nil)
