package action

import (
	"io/ioutil"
	"log"

	"github.com/google/go-cmp/cmp"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"

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

type RulesConfig struct {
	Version int
	Rules   []Rule
}

var rulesConfig RulesConfig

func LoadRulesConfig(filepath string) error {
	log.Printf("Loading rules config file %s", filepath)
	rulesBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	var rawConfig map[string]interface{}
	err = yaml.Unmarshal(rulesBytes, &rawConfig)
	if err != nil {
		return err
	}

	config := &mapstructure.DecoderConfig{
		DecodeHook: actionDecodeHook,
		Result:     &rulesConfig,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	err = decoder.Decode(rawConfig)
	if err != nil {
		return err
	}

	log.Printf("Loaded %d rules", len(rulesConfig.Rules))
	return nil
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
