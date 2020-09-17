package rules

import (
	"log"

	"github.com/google/go-cmp/cmp"

	"github.com/doron-cohen/antidot/internal/dotfile"
)

type Rule struct {
	Name        string
	Description string
	Dotfile     *dotfile.Dotfile
	Ignore      bool
	Actions     []Action
}

// TODO: use some colors
func (r Rule) Pprint() {
	log.Printf("Rule %s:", r.Name)
	for _, action := range r.Actions {
		action.Pprint()
	}

	if r.Ignore {
		log.Println("Rule ignored")
	}
}

func MatchRule(dotfile *dotfile.Dotfile) *Rule {
	for _, rule := range rulesConfig.Rules {
		if cmp.Equal(dotfile, rule.Dotfile) {
			log.Printf("Matched rule %s with dotfile %s", rule.Name, dotfile.Name)
			return &rule
		}
	}
	return nil
}
