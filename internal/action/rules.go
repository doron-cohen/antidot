package action

import (
	"io/ioutil"
	"log"

	"github.com/doron-cohen/antidot/internal/dotfile"
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v2"
)

type Rule struct {
	Name        string
	Description string
	Dotfile     *dotfile.Dotfile
	Ignore      bool
}

type RulesConfig struct {
	Version int
	Rules   []Rule
}

var rulesConfig RulesConfig

func init() {
	filename := "rules.yaml"
	rulesBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read rules file %s: #%v", filename, err)
	}
	err = yaml.Unmarshal(rulesBytes, &rulesConfig)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	log.Printf("Rule %v", rulesConfig)
}

func MatchActions(dotfile *dotfile.Dotfile) {
	for _, rule := range rulesConfig.Rules {
		if cmp.Equal(dotfile, rule.Dotfile) {
			log.Printf("Matched rule %s with dotfile %s", rule, dotfile)
			if rule.Ignore {
				log.Printf("Ignoring dotfile %s", dotfile.Name)
			}
			break
		}
	}
}
