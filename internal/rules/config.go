package rules

import (
	"io/ioutil"
	"os"

	"github.com/doron-cohen/antidot/internal/tui"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type MissingRulesFile struct{}

func (e *MissingRulesFile) Error() string {
	return "Rules file is missing"
}

type RulesConfig struct {
	Version int
	Rules   []Rule
}

var rulesConfig RulesConfig

func LoadRulesConfig(filepath string) (RulesConfig, error) {
	tui.Debug("Loading rules config file %s", filepath)
	rulesBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return RulesConfig{}, &MissingRulesFile{}
		}
		return RulesConfig{}, err
	}

	var rawConfig map[string]interface{}
	err = yaml.Unmarshal(rulesBytes, &rawConfig)
	if err != nil {
		return RulesConfig{}, err
	}

	// We want to completely override the old config
	rulesConfig = RulesConfig{}
	config := &mapstructure.DecoderConfig{
		DecodeHook: actionDecodeHook,
		Result:     &rulesConfig,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return RulesConfig{}, err
	}

	err = decoder.Decode(rawConfig)
	if err != nil {
		return RulesConfig{}, err
	}

	tui.Debug("Loaded %d rules", len(rulesConfig.Rules))
	return rulesConfig, nil
}
