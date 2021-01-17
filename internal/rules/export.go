package rules

import (
	"github.com/doron-cohen/antidot/internal/tui"
)

type Export struct {
	Key   string
	Value string
}

func (e Export) Apply(actx ActionContext) error {
	err := actx.Shell.AddEnvExport(e.Key, e.Value)
	if err != nil {
		return err
	}

	return nil
}

func (e Export) Pprint() {
	tui.Print(
		"  %s %s%s\"%s\"",
		tui.ApplyStyle(tui.Blue, "EXPORT"),
		e.Key,
		tui.ApplyStyle(tui.Gray, "="),
		e.Value,
	)
}

func init() {
	registerAction("export", Export{})
}
