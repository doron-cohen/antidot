package shell

import (
	"fmt"
	"strings"
	"regexp"

	"github.com/doron-cohen/antidot/internal/utils"
)

type Fish struct{}

var bracketedVarRe = regexp.MustCompile(`\$\{(\w+)\}`)

func unbracketEnvVar(str string) string {
	bytes := bracketedVarRe.ReplaceAll([]byte(str), []byte("$$${1}"))
	return string(bytes)
}

func (f *Fish) FormatAlias(alias, command string) string {
	command = unbracketEnvVar(command)
	return fmt.Sprintf("alias %s \"%s\"\n", alias, command)
}

func (f *Fish) FormatExport(key, value string) string {
	value = unbracketEnvVar(value)
	return fmt.Sprintf("set -gx %s \"%s\"\n", key, value)
}

func (f *Fish) InitStub() string {
	format := "set -q %s; or set -x %s \"%s\"\n"

	builder := strings.Builder{}
	builder.WriteString("# Put 'antidot init -s fish | source' (without single quotes) in `fish_config_dir/conf.d/antidot.fish` to automatically run this\n")
	for key, value := range utils.XdgDefaults() {
		builder.WriteString(fmt.Sprintf(format, key, key, value))
	}
	return builder.String()
}

func init() {
	registerShell("fish", &Fish{})
}
