package rules

import (
	"github.com/google/go-cmp/cmp"

	"github.com/doron-cohen/antidot/internal/dotfile"
	"github.com/doron-cohen/antidot/internal/tui"
)

type Rule struct {
	Name        string
	Description string
	Dotfile     *dotfile.Dotfile
	Ignore      bool
	Actions     []Action
	Notes       []string
}

func (r Rule) Pprint() {
	tui.Print(tui.ApplyStylef(tui.Cyan, "Rule %s:", r.Name))
	if len(r.Notes) != 0 {
		for _, note := range r.Notes {
			tui.Print("  %s %s", tui.ApplyStyle(tui.Cyan, "NOTICE"), note)
		}
	}

	for _, action := range r.Actions {
		action.Pprint()
	}

	if r.Ignore {
		tui.Print(tui.ApplyStyle(tui.Gray, "  IGNORED"))
	}
}

func (r Rule) Apply(actx ActionContext) {
	if !r.Ignore {
		for _, action := range r.Actions {
			err := action.Apply(actx)
			if err != nil {
				tui.Warn("Failed to run rule %s: %v", r.Name, err)
				break
			}
		}
	}
}

func MatchRule(dotfile *dotfile.Dotfile) *Rule {
	for _, rule := range rulesConfig.Rules {
		if cmp.Equal(dotfile, rule.Dotfile) {
			tui.Debug("Matched rule %s with dotfile %s", rule.Name, dotfile.Name)
			return &rule
		}
	}
	return nil
}
