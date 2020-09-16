package action

import (
	"log"

	"github.com/doron-cohen/antidot/internal/dotfile"
	"github.com/google/go-cmp/cmp"
)

type Rule struct {
	dotfile *dotfile.Dotfile
	ignore  bool
}

var rules []Rule

func init() {
	rules = make([]Rule, 1)
	rules = append(rules,
		Rule{
			dotfile: dotfile.NewDotfile(".ssh", true),
			ignore:  true,
		},
	)
}

func MatchActions(dotfile *dotfile.Dotfile) {
	for _, rule := range rules {
		if cmp.Equal(dotfile, rule.dotfile) {
			log.Printf("Matched rule %s with dotfile %s", rule, dotfile)
			if rule.ignore {
				log.Printf("Ignoring dotfile %s", dotfile.Name)
			}
			break
		}
	}
}
