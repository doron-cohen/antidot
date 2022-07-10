package shell

import (
	"fmt"
	"regexp"

	"github.com/doron-cohen/antidot/internal/utils"
)

type Fish struct{}

var bracketedVarRe = regexp.MustCompile(`\$\{(\w+)\}`)

func unbracketEnvVar(str string) string {
	bytes := bracketedVarRe.ReplaceAll([]byte(str), []byte("$$${1}"))
	return string(bytes)
}

func (f *Fish) EnvFilePath() (string, error) {
	return utils.AppDirs.GetDataFile("env.fish")
}

func (f *Fish) AliasFilePath() (string, error) {
	return utils.AppDirs.GetDataFile("alias.fish")
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
	envFilePath, _ := f.EnvFilePath()
	aliasFilePath, _ := f.AliasFilePath()

	format := "set -q %s; or set -x %s \"%s\"\n"
	xdgExport := ""
	for key, value := range utils.XdgDefaults() {

		xdgExport += fmt.Sprintf(format, key, key, value)
	}

	return fmt.Sprintf(`%s
if [ -f "%s" ]; source "%s"; end
if [ -f "%s" ]; source "%s"; end`,
		xdgExport,
		envFilePath,
		envFilePath,
		aliasFilePath,
		aliasFilePath,
	)
}

func init() {
	registerShell("fish", &Fish{})
}
