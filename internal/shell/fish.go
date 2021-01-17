package shell

import (
	"fmt"
	"regexp"

	"github.com/doron-cohen/antidot/internal/utils"
)

type Fish struct {
	envMapFormat   *keyValueMapFormat
	aliasMapFormat *keyValueMapFormat
}

var bracketedVarRe = regexp.MustCompile(`\$\{(\w+)\}`)

func unbracketEnvVar(str string) string {
	bytes := bracketedVarRe.ReplaceAll([]byte(str), []byte("$$${1}"))
	return string(bytes)
}

func (f *Fish) envFilePath() (string, error) {
	return utils.AppDirs.GetDataFile("env.fish")
}

func (f *Fish) aliasFilePath() (string, error) {
	return utils.AppDirs.GetDataFile("alias.fish")
}

func (f *Fish) AddEnvExport(key, value string) error {
	path, err := f.envFilePath()
	if err != nil {
		return err
	}

	value = unbracketEnvVar(value)
	return AppendKeyValueToFile(path, key, value, f.envMapFormat)
}

func (f *Fish) AddCommandAlias(alias, command string) error {
	path, err := f.aliasFilePath()
	if err != nil {
		return err
	}

	command = unbracketEnvVar(command)
	return AppendKeyValueToFile(path, alias, command, f.aliasMapFormat)
}

func (f *Fish) InitStub() string {
	envFilePath, _ := f.envFilePath()
	aliasFilePath, _ := f.aliasFilePath()

	format := "set -q %s; or set -x %s=\"%s\"\n"
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
	registerShell("fish", &Fish{
		envMapFormat: NewKeyValueMapFormat(
			`^set -gx (?P<key>\w+) "(?P<value>.*)"`,
			"set -gx %s \"%s\"\n",
		),
		aliasMapFormat: NewKeyValueMapFormat(
			`^alias (?P<alias>\w+) "(?P<command>.*)"`,
			"alias %s \"%s\"\n",
		),
	})
}
