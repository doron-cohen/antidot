package rules

import (
	"github.com/doron-cohen/antidot/internal/tui"
)

type Alias struct {
	Alias   string
	Command string
}

// TODO: remove code duplication with export.go
func (a Alias) Apply(actx ActionContext) error {
	err := actx.Shell.AddCommandAlias(a.Alias, a.Command)
	if err != nil {
		return err
	}

	return nil
}

func (a Alias) Pprint() {
	tui.Print(
		"  %s %s%s\"%s\"",
		tui.ApplyStyle(tui.Magenta, "ALIAS"),
		a.Alias,
		tui.ApplyStyle(tui.Gray, "="),
		a.Command,
	)
}

func init() {
	registerAction("alias", Alias{})
}
